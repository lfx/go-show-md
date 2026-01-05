package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"go-show-md/internal/config"
	"go-show-md/internal/renderer"
)

type ViewHandler struct {
	config *config.Config
	tmpl   *template.Template
}

func NewViewHandler(cfg *config.Config, tmpl *template.Template) *ViewHandler {
	return &ViewHandler{
		config: cfg,
		tmpl:   tmpl,
	}
}

func (h *ViewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		http.Error(w, "Missing file parameter", http.StatusBadRequest)
		return
	}

	if !h.config.IsPathAllowed(filePath) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	html, err := renderer.RenderMarkdown(content)
	if err != nil {
		log.Printf("Error rendering markdown: %v", err)
		http.Error(w, "Failed to render markdown", http.StatusInternalServerError)
		return
	}

	data := struct {
		FilePath string
		Content  template.HTML
		FileName string
	}{
		FilePath: filePath,
		Content:  template.HTML(html),
		FileName: filePath,
	}

	if err := h.tmpl.ExecuteTemplate(w, "view.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
