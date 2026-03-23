package watcher

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"go-show-md/internal/config"
)

type FileEvent struct {
	Path      string
	EventType string
}

type Watcher struct {
	fsWatcher     *fsnotify.Watcher
	cfg           *config.Config
	directories   []string
	eventChan     chan FileEvent
	debounceMap   map[string]*time.Timer
	debounceMutex sync.Mutex
	stopChan      chan bool
}

func New(cfg *config.Config) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		fsWatcher:   fsWatcher,
		cfg:         cfg,
		directories: []string{},
		eventChan:   make(chan FileEvent, 100),
		debounceMap: make(map[string]*time.Timer),
		stopChan:    make(chan bool),
	}, nil
}

func (w *Watcher) AddDirectory(dir string) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	for _, d := range w.directories {
		if d == absDir {
			return nil
		}
	}

	if err := w.addRecursive(absDir, absDir); err != nil {
		return err
	}

	w.directories = append(w.directories, absDir)
	return nil
}

func (w *Watcher) addRecursive(dir string, rootDir string) error {
	globalPatterns := config.GetGlobalIgnorePatterns()
	var localPatterns []string
	if rootDir != "" {
		localPatterns = config.GetLocalIgnorePatterns(rootDir)
	}
	patterns := config.MergeIgnorePatterns(config.DefaultIgnoredPatterns, globalPatterns, localPatterns, w.cfg.IgnoredPatterns)

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if config.IsIgnored(path, info, patterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			if err := w.fsWatcher.Add(path); err != nil {
				log.Printf("Failed to watch directory %s: %v", path, err)
			}
		}

		return nil
	})
}

func (w *Watcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-w.fsWatcher.Events:
				if !ok {
					return
				}

				if strings.HasSuffix(strings.ToLower(event.Name), ".md") {
					w.debounceEvent(event)
				}

				if event.Op&fsnotify.Create != 0 {
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() {
						// Note: for newly created directories inside a watched root, 
						// we'd need to find the actual root to get the local patterns.
						// To keep it simple, we pass the directory itself as root, 
						// which will pick up any local ignore file if it exists inside.
						w.addRecursive(event.Name, event.Name)
					}
				}

			case err, ok := <-w.fsWatcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)

			case <-w.stopChan:
				return
			}
		}
	}()
}

func (w *Watcher) debounceEvent(event fsnotify.Event) {
	w.debounceMutex.Lock()
	defer w.debounceMutex.Unlock()

	if timer, exists := w.debounceMap[event.Name]; exists {
		timer.Stop()
	}

	w.debounceMap[event.Name] = time.AfterFunc(300*time.Millisecond, func() {
		eventType := "modified"
		if event.Op&fsnotify.Create != 0 {
			eventType = "created"
		} else if event.Op&fsnotify.Remove != 0 {
			eventType = "removed"
		} else if event.Op&fsnotify.Rename != 0 {
			eventType = "renamed"
		}

		w.eventChan <- FileEvent{
			Path:      event.Name,
			EventType: eventType,
		}

		w.debounceMutex.Lock()
		delete(w.debounceMap, event.Name)
		w.debounceMutex.Unlock()
	})
}

func (w *Watcher) Events() <-chan FileEvent {
	return w.eventChan
}

func (w *Watcher) Close() error {
	close(w.stopChan)
	return w.fsWatcher.Close()
}
