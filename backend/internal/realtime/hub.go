package realtime

import (
	"log"
	"sync"
)

// Event represents a realtime event that can be pushed to clients.
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Hub is a minimal in-memory hub abstraction. In production, this would be
// backed by Redis pub/sub so that multiple API instances can broadcast to
// all connected clients across the cluster.
type Hub struct {
	mu      sync.Mutex
	clients map[*Client]struct{}
}

type Client struct {
	Send chan Event
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = struct{}{}
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
}

func (h *Hub) Broadcast(ev Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		select {
		case c.Send <- ev:
		default:
			log.Println("dropping event to slow client")
		}
	}
}

