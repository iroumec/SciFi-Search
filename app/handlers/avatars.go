package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Creación del cliente MinIO.
func NewMinioClient() (*s3.Client, error) {
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

// Creación del Bucket y configuración de este como público para lectura de avatares.
func EnsureBucket(ctx context.Context, client *s3.Client, bucket string) error {
	// Crea el bucket si no existe.
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
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
				"Resource": ["arn:aws:s3:::%s/avatars/*"]
			}
		]
	}`, bucket)

	_, err = client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policy),
	})

	return err
}

// Sube el avatar y devuelve la URL pública permanente.
func UploadAvatar(ctx context.Context, client *s3.Client, bucket, userID string, file io.Reader) (string, error) {
	key := "avatars/" + userID + ".jpg"

	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", err
	}

	// URL pública que se guarda en la BD y se usa en el frontend.
	return fmt.Sprintf("%s/%s/%s", os.Getenv("MINIO_HOST"), bucket, key), nil
}

// Se elimina el avatar de un usuario.
func DeleteAvatar(ctx context.Context, client *s3.Client, bucket, userID string) error {
	key := "avatars/" + userID + ".jpg"
	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
