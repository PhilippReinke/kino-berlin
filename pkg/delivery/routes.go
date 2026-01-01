package delivery

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/PhilippReinke/kino-berlin/pkg/app"
)

type Handler struct {
	app       *app.App
	templates *template.Template
	staticDir string
}

func NewHandler(a *app.App, templateDir, staticDir string) (*Handler, error) {
	tmpl, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, err
	}

	return &Handler{
		app:       a,
		templates: tmpl,
		staticDir: staticDir,
	}, nil
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /", http.FileServer(http.Dir(h.staticDir)))
	mux.HandleFunc("GET /api/selects", h.handleSelects)
	mux.HandleFunc("POST /api/screenings", h.handleScreenings)
}
