package notifications

import (
	"encoding/json"
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

func (h *Hub) SendNotification(userID int, message string, tripID int) {
	h.mu.Lock()
	conn, exists := h.clients[userID]
	h.mu.Unlock()

	if exists {
		payload := map[string]interface{}{
			"message": message,
			"trip_id": tripID,
			"type":    "TRIP_READY",
		}
		data, _ := json.Marshal(payload)
		conn.WriteMessage(websocket.TextMessage, data)
	}
}
