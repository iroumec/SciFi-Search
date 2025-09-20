package handlers

import (
	"net/http"
)

// Crea un handler genérico que sirve un template si es GET.
func handlerTemplate(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido.", http.StatusMethodNotAllowed)
			return
		}
		renderizeTemplate(w, path, nil, nil)
	}
}

func registrarHandlersAyuda() {
	rutas := map[string]string{
		"/ayuda":            "template/ayuda/ayuda.html",
		"/albergue":         "template/ayuda/albergue.html",
		"/carnet-deportivo": "template/ayuda/carnet-deportivo.html",
		// Acá se puede seguir agregando.
	}

	for ruta, tmpl := range rutas {
		http.HandleFunc(ruta, handlerTemplate(tmpl))
	}
}
