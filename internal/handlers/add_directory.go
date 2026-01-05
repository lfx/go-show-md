package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go-show-md/internal/config"
	"go-show-md/internal/watcher"
)

type AddDirectoryHandler struct {
	config  *config.Config
	watcher *watcher.Watcher
}

func NewAddDirectoryHandler(cfg *config.Config, w *watcher.Watcher) *AddDirectoryHandler {
	return &AddDirectoryHandler{
		config:  cfg,
		watcher: w,
	}
}

func (h *AddDirectoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Directory string `json:"directory"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Directory == "" {
		http.Error(w, "Directory path is required", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(req.Directory)
	if err != nil {
		http.Error(w, "Directory does not exist", http.StatusBadRequest)
		return
	}

	if !info.IsDir() {
		http.Error(w, "Path is not a directory", http.StatusBadRequest)
		return
	}

	added := h.config.AddDirectory(req.Directory)
	if !added {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Directory already being watched",
		})
		return
	}

	if err := h.config.Save(config.DefaultConfigPath); err != nil {
		log.Printf("Error saving config: %v", err)
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	if err := h.watcher.AddDirectory(req.Directory); err != nil {
		log.Printf("Error adding directory to watcher: %v", err)
		http.Error(w, "Failed to watch directory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Directory added successfully",
	})
}
