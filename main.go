package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"go-show-md/internal/config"
	"go-show-md/internal/handlers"
	"go-show-md/internal/watcher"
)

func main() {
	cfg, err := config.Load(config.DefaultConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	w, err := watcher.New()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}
	defer w.Close()

	for _, dir := range cfg.WatchedDirectories {
		if err := w.AddDirectory(dir); err != nil {
			log.Printf("Failed to watch directory %s: %v", dir, err)
		} else {
			log.Printf("Watching directory: %s", dir)
		}
	}

	w.Start()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/", handlers.NewHomeHandler(cfg, tmpl))
	http.Handle("/view", handlers.NewViewHandler(cfg, tmpl))
	http.Handle("/api/add-directory", handlers.NewAddDirectoryHandler(cfg, w))
	http.Handle("/api/upload", handlers.NewUploadHandler(cfg, w))
	http.Handle("/ws", handlers.NewWSHandler(w))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("Starting server on http://%s", addr)
	log.Printf("Watching %d directories", len(cfg.WatchedDirectories))

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
