package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go-show-md/internal/watcher"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	clients   map[*websocket.Conn]bool
	clientsMu sync.RWMutex
	watcher   *watcher.Watcher
}

func NewWSHandler(w *watcher.Watcher) *WSHandler {
	handler := &WSHandler{
		clients: make(map[*websocket.Conn]bool),
		watcher: w,
	}

	go handler.broadcastFileChanges()

	return handler
}

func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	h.clientsMu.Lock()
	h.clients[conn] = true
	h.clientsMu.Unlock()

	defer func() {
		h.clientsMu.Lock()
		delete(h.clients, conn)
		h.clientsMu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *WSHandler) broadcastFileChanges() {
	for event := range h.watcher.Events() {
		message := map[string]string{
			"type":  "file_changed",
			"path":  event.Path,
			"event": event.EventType,
		}

		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			continue
		}

		h.clientsMu.RLock()
		for client := range h.clients {
			if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				client.Close()
			}
		}
		h.clientsMu.RUnlock()
	}
}
