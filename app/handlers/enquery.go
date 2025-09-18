package handlers

import (
	"net/http"

	gomail "gopkg.in/mail.v2"

	_ "github.com/lib/pq"
)

// ------------------------------------------------------------------------------------------------
// Enquery Handler
// ------------------------------------------------------------------------------------------------

func enqueryHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		enqueryHandleGET(w, "")
	case http.MethodPost:
		enqueryHandlePOST(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func enqueryHandleGET(w http.ResponseWriter, errorMessage string) {

	data := map[string]any{
		"ErrorMessage": errorMessage,
	}

	renderizeTemplate(w, "template/enquery/enquery.html", data)
}

// ------------------------------------------------------------------------------------------------

func enqueryHandlePOST(w http.ResponseWriter, r *http.Request) {

	// Se parsean los datos del formulario enviados vía POST.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	surname := r.FormValue("surname")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	address := r.FormValue("address")
	enquery := r.FormValue("enquery")

	if isThereEmptyField(email, enquery) {
		enqueryHandleGET(w, "Faltan campos obligatorios.")
		return
	}

	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "consultas@olimpiadas.com")
	message.SetHeader("To", "iroumec@alumnos.exa.unicen.edu.ar")
	message.SetHeader("Subject", "Consulta - Olimpiadas")

	// Set email body
	emailBody := `
		Consulta automatizada de la página de la facultad.
	`
	message.SetBody("text/plain", "This is the Test Body")

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", "1a2b3c4d5e6f7g")

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		enqueryHandleGET(w, "Ha ocurrido un erro al enviar la consulta. Tranquilo. ¡La culpa no es tuya! Intenta envviar el email directamente.")
		panic(err)
		return
	}

	// El email fue enviado exitosamente.

	renderizeTemplate(w, "template/enquery/enquerySent.html", nil)
}
