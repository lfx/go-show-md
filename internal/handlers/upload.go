package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go-show-md/internal/config"
	"go-show-md/internal/watcher"
)

type UploadHandler struct {
	config  *config.Config
	watcher *watcher.Watcher
}

func NewUploadHandler(cfg *config.Config, w *watcher.Watcher) *UploadHandler {
	return &UploadHandler{
		config:  cfg,
		watcher: w,
	}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	watchedDir := "./watched-files"
	if err := os.MkdirAll(watchedDir, 0755); err != nil {
		log.Printf("Error creating watched-files directory: %v", err)
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	added := h.config.AddDirectory(watchedDir)
	if added {
		if err := h.config.Save(config.DefaultConfigPath); err != nil {
			log.Printf("Error saving config: %v", err)
		}

		if err := h.watcher.AddDirectory(watchedDir); err != nil {
			log.Printf("Error adding directory to watcher: %v", err)
		}
	}

	destPath := filepath.Join(watchedDir, header.Filename)
	dest, err := os.Create(destPath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		log.Printf("Error copying file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	absPath, _ := filepath.Abs(destPath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "File uploaded successfully",
		"path":    absPath,
	})
}
