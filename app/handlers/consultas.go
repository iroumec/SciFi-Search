package handlers

import (
	"fmt"
	"net/http"

	gomail "gopkg.in/mail.v2"

	_ "github.com/lib/pq"
)

// ------------------------------------------------------------------------------------------------
// Registro de Handlers de Consultas
// ------------------------------------------------------------------------------------------------

func registrarHandlersConsultas() {

	http.HandleFunc("/consultar", manejarConsultas)
}

// ------------------------------------------------------------------------------------------------
// Manejo de Consultas
// ------------------------------------------------------------------------------------------------

func manejarConsultas(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		mostrarFormularioConsulta(w, "")
	case http.MethodPost:
		procesarConsulta(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func mostrarFormularioConsulta(w http.ResponseWriter, errorMessage string) {

	data := map[string]any{
		"ErrorMessage": errorMessage,
	}

	renderizeTemplate(w, "template/consultas/consulta.html", data, nil)
}

// ------------------------------------------------------------------------------------------------

func procesarConsulta(w http.ResponseWriter, r *http.Request) {

	// Se parsean los datos del formulario enviados vía POST.
	if err := r.ParseForm(); err != nil {
		mostrarFormularioConsulta(w, "Error al procesar el formulario.")
		return
	}

	// Se obtienen los datos del formulatio.
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	address := r.FormValue("address")
	enquery := r.FormValue("enquery")

	// Se verifica que no haya campos obligatorios vacíos.
	if hayCampoIncompleto(email, enquery) {
		mostrarFormularioConsulta(w, "Faltan campos obligatorios.")
		return
	}

	// Se crea un nuevo mensaje.
	message := gomail.NewMessage()

	// Se setean los encabezados del email.
	message.SetHeader("From", "consultas@olimpiadas.com")
	message.SetHeader("To", "iroumec@alumnos.exa.unicen.edu.ar")
	message.SetHeader("Subject", "Consulta - Olimpiadas")

	// Se define el cuerpo del mensaje.
	emailBodyFormat := `
		Consulta automatizada de la página de la facultad.

		Nombre: %s
		Apellido: %s
		Email: %s
		Teléfono: %s
		Dirección: %s

		Consulta:
		
		%s
	`

	// Se le da formato al cuerpo del mensaje.
	emailBody := fmt.Sprintf(emailBodyFormat, name, surname, email, phone, address, enquery)

	// Se establece el cuerpo del mensaje.
	message.SetBody("text/plain", emailBody)

	/*

		TODO: arreglar para poder enviar luego un email directamente.

		// Set up the SMTP dialer
			dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", "1a2b3c4d5e6f7g")

			// Send the email
			if err := dialer.DialAndSend(message); err != nil {
				enqueryHandleGET(w, "Ha ocurrido un error al enviar la consulta. Tranquilo. ¡La culpa no es tuya! Intenta envviar el email directamente.")
				panic(err)
				return
			}

	*/

	// El email fue enviado exitosamente.
	renderizeTemplate(w, "template/consultas/consulta-enviada.html", nil, nil)
}
