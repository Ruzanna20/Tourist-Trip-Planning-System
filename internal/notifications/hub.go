package notifications

import (
	"fmt"
	"net/http"
	"sync"
)

type Hub struct {
	clients map[int]chan []byte
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int]chan []byte),
	}
}

func (h *Hub) HandleSSE(w http.ResponseWriter, r *http.Request, userID int) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan []byte)

	h.mu.Lock()
	h.clients[userID] = messageChan
	h.mu.Unlock()

	notify := r.Context().Done()

	go func() {
		<-notify
		h.mu.Lock()
		delete(h.clients, userID)
		close(messageChan)
		h.mu.Unlock()
	}()

	for {
		msg, ok := <-messageChan
		if !ok {
			break
		}
		fmt.Fprintf(w, "data: %s\n\n", msg)
		fmt.Fprintf(w, ": heartbeat\n\n")
		w.(http.Flusher).Flush()
	}
}

func (h *Hub) SendNotification(userID int, data []byte) {
	h.mu.Lock()
	ch, exists := h.clients[userID]
	h.mu.Unlock()

	if exists {
		ch <- data
	}
}
