package handlers

import (
	"net/http"
)

// Handler genérico para GET
func handlerTemplate(path string, data map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido.", http.StatusMethodNotAllowed)
			return
		}
		renderizeTemplate(w, path, data, nil)
	}
}

func registrarHandlersAyuda() {
	rutas := map[string]struct {
		tmpl string
		data map[string]any
	}{
		"/ayuda": {
			tmpl: "template/ayuda/ayuda.html",
		},
		"/albergue": {
			tmpl: "template/ayuda/albergue.html",
			data: map[string]any{
				"fotos": obtenerFotos("static/img/albergue/"),
			},
		},
		"/carnet-deportivo": {
			tmpl: "template/ayuda/carnet-deportivo.html",
		},
	}

	for ruta, def := range rutas {
		http.HandleFunc(ruta, handlerTemplate(def.tmpl, def.data))
	}
}
