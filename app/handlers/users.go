package handlers

// ------------------------------------------------------------------------------------------------

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	sqlc "tpe/web/app/database"
	"tpe/web/app/utils"
	"tpe/web/app/views"

	"github.com/a-h/templ"
)

// ------------------------------------------------------------------------------------------------

// Se registran los endpoints relacionados al manejo de usarios.
func registrarHandlersUsuarios() {

	http.HandleFunc("/users", userHandler)
}

// ------------------------------------------------------------------------------------------------

func userHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		if utils.HasGETRequestParameters(r) {
			showUser(w, r)
		} else {
			listUsers(w, r)
		}
	case http.MethodPost:
		addUser(w, r)
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

// Agrega un usuario a la base de datos.
func addUser(w http.ResponseWriter, r *http.Request) {

	// Se parsea el body de la request.
	err := r.ParseForm()
	if err != nil {
		// Cuerpo mal formado
	}

	name := r.FormValue("name")
	middlename := r.FormValue("middlename")
	surname := r.FormValue("surname")

	if hayCampoIncompleto(name, surname) {
		// Campos obligatorios incompletos -> 400 Bad Request
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}
	_, err = queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Name:       name,
		Middlename: sql.NullString{String: middlename, Valid: middlename != ""},
		Surname:    surname,
	})
	if err != nil {
		// Error al crear el usuario en la BD -> 500 Internal Server Error.
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Se establece el código de estado a 201 Created.
	w.WriteHeader(http.StatusCreated)

	// Acá tendría que haber un "usuario registrado con éxito".
	component := views.SuccessfulSignUpPage()
	templ.Handler(component).ServeHTTP(w, r)
}

// ------------------------------------------------------------------------------------------------

// Elimina un usuario de la base de datos.
func deleteUser(w http.ResponseWriter, r *http.Request) {

	id, err := extractID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = queries.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: El usuario no existe.
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			// Error 500: Hubo un problema con la base de datos u otro error inesperado.
			log.Printf("Error al obtener usuario por ID %d: %v", id, err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	component := views.UserDeletedPage()
	templ.Handler(component).ServeHTTP(w, r)
}

// ------------------------------------------------------------------------------------------------

// Lista a todos los usuarios.
func listUsers(w http.ResponseWriter, r *http.Request) {

	users, err := queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := views.UserPage(users)
	templ.Handler(component).ServeHTTP(w, r)
}

// ------------------------------------------------------------------------------------------------

// Muestra los datos correspondientes a un usuario, dado un ID.
func showUser(w http.ResponseWriter, r *http.Request) {

	id, err := extractID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := queries.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: El usuario no existe.
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			// Error 500: Hubo un problema con la base de datos u otro error inesperado.
			log.Printf("Error al obtener usuario por ID %d: %v", id, err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	component := views.ProfilePage(user)
	templ.Handler(component).ServeHTTP(w, r)
}

// ------------------------------------------------------------------------------------------------

func updateUser(w http.ResponseWriter, r *http.Request) {

	id, err := extractID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

// ------------------------------------------------------------------------------------------------

func extractID(r *http.Request) (int32, error) {

	// Obtención del valor del parámetro 'id' directamente.
	idString := r.URL.Query().Get("id")
	if idString == "" {
		return 0, fmt.Errorf("parámetro 'id' es requerido")
	}

	// Conversión del ID de string a un número, validando que quepa en 32 bits.
	id64, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parámetro 'id' debe ser un número entero válido")
	}

	// Si todo fue exitoso, se convierte el id a int32.
	return int32(id64), nil
}

// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------
// SignIn Handler
// ------------------------------------------------------------------------------------------------

/*
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

	// TODO: podría limitarse el tamaño del archivo en lugar de cuánto se gaurda en memoria.
	// Los certificados siempre van a pesar poco.

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

	// Se encripta la contraseña para no manejar credenciales en bruto.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	_, err = queries.CrearUsuario(r.Context(), sqlc.CrearUsuarioParams{
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

	// Se compara la contraseña con la almacenada en el servidor.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Contraseña), []byte(password)); err != nil {
		mostrarFormularioLogin(w, "Contraseña incorrecta.")
		return
	}

	handleProfileAccess(user, w, r)
}
*/
