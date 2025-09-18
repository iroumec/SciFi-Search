package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

// ------------------------------------------------------------------------------------------------
// Review Handler
// ------------------------------------------------------------------------------------------------

func handlerComentarios(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("template/profile.html"))

	fmt.Println("\nManejando renderizado del perfil...")

	c, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "no autenticado", http.StatusUnauthorized)
		return
	}

	userID, ok := sessions[c.Value]
	if !ok {
		http.Error(w, "sesión inválida", http.StatusUnauthorized)
		return
	}

	user, err := queries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "error cargando usuario", http.StatusInternalServerError)
		return
	}

	fmt.Println("Renderizando plantilla de acuerdo a los datos del usuario...")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
	}

	fmt.Println("Plantilla renderizada.")
}
