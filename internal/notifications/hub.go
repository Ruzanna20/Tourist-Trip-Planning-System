package notifications

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients map[int]*websocket.Conn
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int]*websocket.Conn),
	}
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request, userID int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.mu.Lock()
	h.clients[userID] = conn
	h.mu.Unlock()

}

func (h *Hub) SendNotification(userID int, data []byte) {
	h.mu.Lock()
	conn, exists := h.clients[userID]
	h.mu.Unlock()

	if exists {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			slog.Error("Failed to send WS message", "userID", userID, "error", err)
			conn.Close()
			h.mu.Lock()
			delete(h.clients, userID)
			h.mu.Unlock()
		}
	} else {
		slog.Warn("WS: User not connected", "userID", userID)
	}
}
