package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins
		},
	}
)

// Message represents a WebSocket message
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// TimerMessage represents a timer update message
type TimerMessage struct {
	Type      string    `json:"type"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"`
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	GameID string
	Send   chan []byte
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToGame sends a message to all clients in a specific game
func (h *Hub) BroadcastToGame(gameID string, messageType string, payload []byte) {
	message := Message{
		Type:    messageType,
		Payload: payload,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.mu.RLock()
	for client := range h.clients {
		if client.GameID == gameID {
			select {
			case client.Send <- messageBytes:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
	h.mu.RUnlock()
}

// GetGameClients returns the number of connected clients for a game
func (h *Hub) GetGameClients(gameID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for client := range h.clients {
		if client.GameID == gameID {
			count++
		}
	}
	return count
}

// CloseGame closes all connections for a game
func (h *Hub) CloseGame(gameID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		if client.GameID == gameID {
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// StartTimer starts a timer for a game
func (h *Hub) StartTimer(gameID string, timerType string, duration int) {
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Second)

	message := TimerMessage{
		Type:      timerType,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling timer message: %v", err)
		return
	}

	h.BroadcastToGame(gameID, "timer_started", messageBytes)

	// Start a goroutine to send timer updates
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				remaining := int(endTime.Sub(time.Now()).Seconds())
				if remaining <= 0 {
					h.BroadcastToGame(gameID, "timer_ended", messageBytes)
					return
				}

				update := TimerMessage{
					Type:      timerType,
					StartTime: startTime,
					EndTime:   endTime,
					Duration:  remaining,
				}

				updateBytes, err := json.Marshal(update)
				if err != nil {
					log.Printf("Error marshaling timer update: %v", err)
					return
				}

				h.BroadcastToGame(gameID, "timer_update", updateBytes)
			}
		}
	}()
}

// Register registers a new client with the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Handler handles WebSocket connections
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// ServeWS handles WebSocket requests from clients
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	if gameID == "" {
		http.Error(w, "game_id is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		GameID: gameID,
	}

	client.Hub.register <- client

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Hub.broadcast <- message
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)
