package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go-show-md/internal/config"
	"go-show-md/internal/watcher"
)

type RemoveDirectoryHandler struct {
	config  *config.Config
	watcher *watcher.Watcher
}

func NewRemoveDirectoryHandler(cfg *config.Config, w *watcher.Watcher) *RemoveDirectoryHandler {
	return &RemoveDirectoryHandler{
		config:  cfg,
		watcher: w,
	}
}

func (h *RemoveDirectoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Remove from config
	removed := h.config.RemoveDirectory(req.Directory)
	if !removed {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Directory is not being watched",
		})
		return
	}

	if err := h.config.Save(config.DefaultConfigPath); err != nil {
		log.Printf("Error saving config after removal: %v", err)
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	// Remove from fsnotify watcher
	if err := h.watcher.RemoveDirectory(req.Directory); err != nil {
		log.Printf("Error removing directory from watcher: %v", err)
		// We still return success as it's logically unwatched
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Directory untracked successfully",
	})
}
