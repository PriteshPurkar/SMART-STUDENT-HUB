package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/example/scalable-learning-platform/backend/internal/realtime"
)

var defaultHub = realtime.NewHub()

func registerRealtimeRoutes(r chi.Router) {
	r.Get("/stream", handleRealtimeStream)
}

// handleRealtimeStream uses a very simple Server-Sent Events style stream
// for the scaffold. In production you could upgrade this to WebSockets and
// back it with Redis pub/sub for cross-instance fan-out.
func handleRealtimeStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	client := &realtime.Client{
		Send: make(chan realtime.Event, 8),
	}
	defaultHub.Register(client)
	defer defaultHub.Unregister(client)

	enc := json.NewEncoder(w)
	for {
		select {
		case ev := <-client.Send:
			// SSE format: data: <json>\n\n
			w.Write([]byte("data: "))
			if err := enc.Encode(ev); err != nil {
				return
			}
			w.Write([]byte("\n"))
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

