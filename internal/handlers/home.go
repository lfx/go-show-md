package handlers

import (
	"html/template"
	"log"
	"net/http"

	"go-show-md/internal/config"
)

type HomeHandler struct {
	config *config.Config
	tmpl   *template.Template
}

func NewHomeHandler(cfg *config.Config, tmpl *template.Template) *HomeHandler {
	return &HomeHandler{
		config: cfg,
		tmpl:   tmpl,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files, err := h.config.ScanMarkdownFiles()
	if err != nil {
		log.Printf("Error scanning files: %v", err)
		http.Error(w, "Failed to scan files", http.StatusInternalServerError)
		return
	}

	filesByDir := make(map[string][]config.FileInfo)
	for _, file := range files {
		filesByDir[file.Directory] = append(filesByDir[file.Directory], file)
	}

	data := struct {
		Directories        []string
		FilesByDirectory   map[string][]config.FileInfo
		WatchedDirectories []string
	}{
		Directories:        getSortedKeys(filesByDir),
		FilesByDirectory:   filesByDir,
		WatchedDirectories: h.config.WatchedDirectories,
	}

	if err := h.tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func getSortedKeys(m map[string][]config.FileInfo) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
