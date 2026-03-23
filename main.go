package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"go-show-md/internal/config"
	"go-show-md/internal/handlers"
	"go-show-md/internal/watcher"

	"github.com/getlantern/systray"
)

// Global state to manage server lifecycle
var (
	server       *http.Server
	serverCtx    context.Context
	serverCancel context.CancelFunc
	serverWg     sync.WaitGroup
	mu           sync.Mutex
	isRunning    bool

	// Configuration and dependencies
	cfg  *config.Config
	tmpl *template.Template
)

func main() {
	// Ensure we run from the executable directory (for .app bundles)
	if exe, err := os.Executable(); err == nil {
		if err := os.Chdir(filepath.Dir(exe)); err != nil {
			log.Printf("Failed to change working directory: %v", err)
		}
	}

	// 1. Load configuration and templates immediately
	var err error
	cfg, err = config.Load(config.DefaultConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// 2. Run the systray application
	// This takes over the main thread
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIconData())
	systray.SetTitle("MD Viewer")
	systray.SetTooltip("Go Markdown Viewer")

	mOpen := systray.AddMenuItem("Open in Browser", "Open the viewer in your default browser")
	mToggle := systray.AddMenuItem("Stop Server", "Turn the server on or off")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Start the server initially
	startServer()
	mToggle.SetTitle("Stop Server")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				openBrowser(fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port))
			case <-mToggle.ClickedCh:
				// We need to check state and toggle
				// We unlock inside the check to avoid holding lock during start/stop operations if they take time
				mu.Lock()
				running := isRunning
				mu.Unlock()

				if running {
					stopServer()
					mToggle.SetTitle("Start Server")
				} else {
					startServer()
					mToggle.SetTitle("Stop Server")
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	stopServer()
}

func startServer() {
	mu.Lock()
	defer mu.Unlock()

	if isRunning {
		return
	}

	// Create a new context for this server run
	serverCtx, serverCancel = context.WithCancel(context.Background())

	// Initialize watcher
	w, err := watcher.New(cfg)
	if err != nil {
		log.Printf("Failed to create watcher: %v", err)
		return
	}

	for _, dir := range cfg.WatchedDirectories {
		if err := w.AddDirectory(dir); err != nil {
			log.Printf("Failed to watch directory %s: %v", dir, err)
		} else {
			log.Printf("Watching directory: %s", dir)
		}
	}
	w.Start()

	// Setup Handler with a NEW ServeMux to avoid collisions on restart
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/", handlers.NewHomeHandler(cfg, tmpl))
	mux.Handle("/view", handlers.NewViewHandler(cfg, tmpl))
	mux.Handle("/api/add-directory", handlers.NewAddDirectoryHandler(cfg, w))
	mux.Handle("/api/upload", handlers.NewUploadHandler(cfg, w))
	mux.Handle("/ws", handlers.NewWSHandler(w))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server = &http.Server{
		Addr:    addr,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return serverCtx
		},
	}

	log.Printf("Starting server on http://%s", addr)
	log.Printf("Watching %d directories", len(cfg.WatchedDirectories))

	serverWg.Add(1)
	go func() {
		defer serverWg.Done()
		defer w.Close() // Cleanup watcher when server stops

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server failed: %v", err)
		}
	}()

	isRunning = true
}

func stopServer() {
	mu.Lock()
	defer mu.Unlock()

	if !isRunning {
		return
	}

	log.Println("Stopping server...")

	// Create a context with timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	if serverCancel != nil {
		serverCancel()
	}

	isRunning = false
	log.Println("Server stopped")
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}
