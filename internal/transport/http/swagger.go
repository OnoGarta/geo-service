package http

import (
	"embed"
	_ "embed"
	"io/fs"
	"net/http"
)

//go:embed swagger.yaml
var specYAML []byte

//go:embed swagger_index.html
var indexHTML string

//go:embed swaggerui/*
var uiFS embed.FS // файлы лежат как ui/swagger-ui.css …

func SwaggerRouter() http.Handler {
	mux := http.NewServeMux()

	// YAML
	mux.HandleFunc("/swagger/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(specYAML)
	})

	subFS, _ := fs.Sub(uiFS, "swaggerui")

	mux.Handle("/swagger/static/",
		http.StripPrefix("/swagger/static/",
			http.FileServer(http.FS(subFS))),
	)

	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(indexHTML))
	})

	return mux
}
