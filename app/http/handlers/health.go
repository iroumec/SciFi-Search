package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import "net/http"

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

func RegisterHealth() {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		// Solo se responde a peticiones GET.
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Se establece el código de estado 200 OK.
		// A esto lo busca `curl -f` cuando se levanta el servidor.
		w.WriteHeader(http.StatusOK)

		// Cuerpo simple para saber que funciona si se abre desde un navegador.
		w.Write([]byte("Servidor OK"))
	})
}

// ------------------------------------------------------------------------------------------------
