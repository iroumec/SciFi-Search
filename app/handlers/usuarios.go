package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"golang.org/x/crypto/bcrypt"

	utils "tpe/web/app/utils"

	sqlc "tpe/web/app/database"

	_ "github.com/lib/pq"
)

/*
Se registran todos los endpoints relacionados al
registro e inicio de sesión de usuarios.
*/
func registerUserHandlers() {

	// Handler que maneja el registro de usuarios.
	http.HandleFunc("/registrarse", registrarUsuario)

	// Handler que maneja el login de usuarios.
	http.HandleFunc("/iniciar-sesion", iniciarSesion)
}

// ------------------------------------------------------------------------------------------------
// SignIn Handler
// ------------------------------------------------------------------------------------------------

func registrarUsuario(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		mostrarFormularioRegistro(w, "")
	case http.MethodPost:
		procesarRegistro(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func mostrarFormularioRegistro(w http.ResponseWriter, errorMessage string) {

	data := map[string]any{
		"ErrorMessage": errorMessage,
	}

	renderizeTemplate(w, "template/usuarios/registro/registrarse.html", data, nil)
}

// ------------------------------------------------------------------------------------------------

func procesarRegistro(w http.ResponseWriter, r *http.Request) {

	// Se parsean los datos del formulario enviados vía POST.
	// 10 es el número base y 20 es el desplazamiento en bit a izquierda.
	// 10 * 2^20 = 10 * 1.048.576 = 10.485.760 bytes = 10 MB.
	// De esta forma, se limita el tamaño del PDF que se suba a 10 MB para no saturar la memoria.
	// Si el formulario es más grande, el resto se guarda automáticamente
	// en archivos temporales en disco (r.MultipartForm.File).
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// TODO: podría limitarse el tamaño del archivo.

	// DNI, nombre y apellido no deberían pedirse. Se obtienen del certificado. TODO
	dni := r.FormValue("dni")
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Se verifica que ninguno de los campos esté incompleto.
	if hayCampoIncompleto(dni, name, email, password) {
		mostrarFormularioRegistro(w, "Faltan campos obligatorios.")
		return
	}

	// Se lee tanto la parte en memoria como la parte en disco.
	// Se obtiene el archivo del formulario.
	file, _, err := r.FormFile("certificado")
	if err != nil {
		// De no haber adjuntado un certificado de alumno regular,
		// se muestra un mensaje de error.
		mostrarFormularioRegistro(w, "Debe adjuntar el certificado de alumno regular.")
		return
	}

	// Independientemente de lo que pase, se cierra el archivo.
	defer file.Close()

	// Validación del PDF.
	valido, err := utils.ValidarConstancia(file)
	if err != nil {
		mostrarFormularioRegistro(w, err.Error())
		return
	}
	if !valido {
		mostrarFormularioRegistro(w, "El certificado no es válido.")
		return
	}

	// Se encripta la contraseña para no manejar credenciales en bruto.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	createdUser, err := queries.CrearUsuario(r.Context(), sqlc.CrearUsuarioParams{
		Dni:        dni,
		Nombre:     name,
		Email:      email,
		Contraseña: string(hashedPassword),
	})
	if err != nil { // Esto anda MUY MAL. TODO: solucionar luego.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			http.Error(w, "El usuario o email ya existen.", http.StatusConflict)
			return
		}

		log.Printf("error creando usuario: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	// Quizás pueda usarse después...
	_ = createdUser

	renderizeTemplate(w, "template/usuarios/registro/registro-exitoso.html", nil, nil)
}

// ------------------------------------------------------------------------------------------------
// LogIn Handler
// ------------------------------------------------------------------------------------------------

func iniciarSesion(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		mostrarFormularioLogin(w, "")
	case http.MethodPost:
		procesarLogin(w, r)
	default:
		http.Error(w, "Método no permitido.", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func mostrarFormularioLogin(w http.ResponseWriter, errorMessage string) {

	data := map[string]any{
		"ErrorMessage": errorMessage,
	}

	renderizeTemplate(w, "template/usuarios/iniciar-sesion.html", data, nil)
}

// ------------------------------------------------------------------------------------------------

func procesarLogin(w http.ResponseWriter, r *http.Request) {

	// Se parsean los datos del formulario enviados vía POST.
	if err := r.ParseForm(); err != nil {
		// Se renderiza la página con el error correspondiente.
		mostrarFormularioLogin(w, "Error al procesar el formulario.")
		return
	}

	dni := r.FormValue("dni")
	password := r.FormValue("password")

	if hayCampoIncompleto(dni, password) {
		mostrarFormularioLogin(w, "Faltan campos obligatorios.")
		return
	}

	// Hacer funcionar esto...
	user, err := queries.ObtenerUsuarioPorDNI(r.Context(), dni)
	if err != nil {
		if err == sql.ErrNoRows {
			mostrarFormularioLogin(w, "El usuario proporcionado no existe.")
			return
		}
		log.Printf("error getting user: %v", err)
		mostrarFormularioLogin(w, "Error interno del servidor.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Contraseña), []byte(password)); err != nil {
		mostrarFormularioLogin(w, "Contraseña incorrecta.")
		return
	}

	handleProfileAccess(user, w, r)
}
