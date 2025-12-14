package handlers

// ------------------------------------------------------------------------------------------------

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"scifi-search/app/database"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
)

// ------------------------------------------------------------------------------------------------

const (
	bucketName = "avatars"
)

// ------------------------------------------------------------------------------------------------

var S3Client *s3.Client

// ------------------------------------------------------------------------------------------------

func registerAvatarHandlers() {

	// Inicialización del cliente.
	client, err := newMinioClient()
	if err != nil {
		return
	}
	S3Client = client

	// Creación del bucket.
	err = ensureBucket(context.Background(), bucketName)
	if err != nil {
		fmt.Println("Error creando bucket:", err)
	}

	http.HandleFunc("/avatar", avatarHandler)
}

// ------------------------------------------------------------------------------------------------

func avatarHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		uploadAvatarHandler(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

// Creación del cliente MinIO.
func newMinioClient() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				os.Getenv("MINIO_ROOT_USER"),
				os.Getenv("MINIO_ROOT_PASSWORD"),
				"",
			),
		),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("MINIO_HOST"))
		o.UsePathStyle = true
	})

	return client, nil
}

// ------------------------------------------------------------------------------------------------

// Creación del bucket y configuración de este como público para lectura de avatares.
func ensureBucket(ctx context.Context, bucket string) error {

	// Se crea el bucket si no existe.
	_, err := S3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil && !strings.Contains(err.Error(), "BucketAlreadyOwned") &&
		!strings.Contains(err.Error(), "BucketAlreadyExists") {
		return err
	}

	// Política pública: cualquiera puede ver los avatares.
	// Pero solo los usuarios autenticados pueden subir.
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": "*"},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucket)

	_, err = S3Client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policy),
	})

	return err
}

// ------------------------------------------------------------------------------------------------

func uploadAvatarHandler(w http.ResponseWriter, r *http.Request) {

	// Parseo del formulario.
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error procesando el archivo", 400)
		return
	}

	// Recuperación del archivo.
	file, _, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "No se pudo leer el archivo", 400)
		return
	}
	defer file.Close()

	// Obtención del usuario.
	user := getCurrentUser(w, r)
	if user == nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		// Acá podría enviarse a la pestaña de login si el usuario no está autenticado.
		// Pero ¿puede el usuario acceder a esto si no está autenticado?
		return
	}

	// Se cambia el tamaño de la imagen.
	resizedFile, err := ResizeImageToAvatar(file)
	if err != nil {
		http.Error(w, "Error procesando imagen", 500)
		return
	}

	// Se sube el archivo al almacenamiento de objetos.
	url, err := UploadAvatar(r.Context(), bucketName, user.UserID, resizedFile)
	if err != nil {
		http.Error(w, "Error subiendo avatar", 500)
		log.Printf("%s", err)
		return
	}

	// Guardado de la URL en la Base de Datos.
	err = queries.UploadAvatar(r.Context(), database.UploadAvatarParams{
		UserID: user.UserID,
		AvatarUrl: sql.NullString{
			String: url,
			Valid:  true,
		},
	})

	// Redirección del usuario a su perfil.
	// TODO: usar HTMX para cargar solo la parte que cambió.
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func UploadAvatar(ctx context.Context, bucket string, userID int32, file io.Reader) (string, error) {

	// Conversión en bytes para obtener un ReadSeeker.
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	body := bytes.NewReader(data) // ReadSeeker. Necesario ya que no se usa TLS.

	key := fmt.Sprintf("%d.jpg", userID)

	_, err = S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", os.Getenv("MINIO_PUBLIC_URL"), bucket, key), nil
}

// ------------------------------------------------------------------------------------------------

// Se elimina el avatar de un usuario.
func deleteAvatar(ctx context.Context, userID int32) error {
	key := fmt.Sprintf("%d.jpg", userID)
	_, err := S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	return err
}

// ------------------------------------------------------------------------------------------------

func ResizeImageToAvatar(file io.Reader) (io.Reader, error) {

	// Decode de la imagen.
	img, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}

	// Resize.
	resized := imaging.Resize(img, 256, 256, imaging.Lanczos)

	// Codificación en JPEG.
	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, resized, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
