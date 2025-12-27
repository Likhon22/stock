package websocket

import "log"

type Hub struct {
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	log.Println("hub started")
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.Printf("Client registered. Total: %d", len(h.clients))
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client unregistered. Total: %d", len(h.clients))
			}
		case message := <-h.Broadcast:
			// Send to all clients
			for client := range h.clients {
				select {
				case client.send <- message:
					// Message sent successfully
				default:
					// Client's buffer is full, disconnect them
					close(client.send)
					delete(h.clients, client)
				}
			}

		}
	}

}
