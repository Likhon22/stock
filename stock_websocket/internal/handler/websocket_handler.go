package handler

import (
	"log"
	"net/http"

	"stock_websocket/internal/websocket"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub *websocket.Hub
}

func NewHandler(hub *websocket.Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Create client
	client := websocket.NewClient(h.hub, conn)

	// Register with hub
	h.hub.Register <- client

	// Start goroutines
	go client.WritePump()
	go client.ReadPump()
}
