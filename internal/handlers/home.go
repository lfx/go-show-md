package handlers

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

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

	sortBy := r.URL.Query().Get("sort")
	if sortBy == "updated" {
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModifiedAt.After(files[j].ModifiedAt)
		})
	} else {
		sortBy = "name"
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name < files[j].Name
		})
	}

	filesByDir := make(map[string][]config.FileInfo)
	for _, file := range files {
		filesByDir[file.Directory] = append(filesByDir[file.Directory], file)
	}

	data := struct {
		Directories        []string
		FilesByDirectory   map[string][]config.FileInfo
		WatchedDirectories []string
		CurrentSort        string
	}{
		Directories:        getSortedKeys(filesByDir, sortBy),
		FilesByDirectory:   filesByDir,
		WatchedDirectories: h.config.WatchedDirectories,
		CurrentSort:        sortBy,
	}

	if err := h.tmpl.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func getSortedKeys(m map[string][]config.FileInfo, sortBy string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	
	if sortBy == "updated" {
		sort.Slice(keys, func(i, j int) bool {
			var timeI, timeJ time.Time
			if len(m[keys[i]]) > 0 {
				timeI = m[keys[i]][0].ModifiedAt
			}
			if len(m[keys[j]]) > 0 {
				timeJ = m[keys[j]][0].ModifiedAt
			}
			if timeI.Equal(timeJ) {
				return keys[i] < keys[j]
			}
			return timeI.After(timeJ)
		})
	} else {
		sort.Strings(keys)
	}
	
	return keys
}
