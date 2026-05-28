package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	coderws "github.com/coder/websocket"
)

const (
	writeWait      = 10 * time.Second
	pingPeriod     = 30 * time.Second
	maxMessageSize = 4 << 10
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type TimerMessage struct {
	Type      string    `json:"type"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"`
}

type Client struct {
	Hub    *Hub
	Conn   *coderws.Conn
	GameID string
	Send   chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

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
		}
	}
}

func (h *Hub) BroadcastToGame(gameID string, messageType string, payload []byte) {
	message := Message{Type: messageType, Payload: payload}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		slog.Error("failed to marshal websocket message", "error", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.GameID != gameID {
			continue
		}
		select {
		case client.Send <- messageBytes:
		default:
			go func(c *Client) { h.unregister <- c }(client)
		}
	}
}

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

func (h *Hub) CloseGame(gameID string) {
	h.mu.RLock()
	var clients []*Client
	for client := range h.clients {
		if client.GameID == gameID {
			clients = append(clients, client)
		}
	}
	h.mu.RUnlock()
	for _, client := range clients {
		h.unregister <- client
	}
}

func (h *Hub) StartTimer(gameID string, timerType string, duration int) {
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(duration) * time.Second)
	message := TimerMessage{Type: timerType, StartTime: startTime, EndTime: endTime, Duration: duration}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		slog.Error("failed to marshal timer message", "error", err)
		return
	}
	h.BroadcastToGame(gameID, "timer_started", messageBytes)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			remaining := int(time.Until(endTime).Seconds())
			if remaining <= 0 {
				h.BroadcastToGame(gameID, "timer_ended", messageBytes)
				return
			}
			update := TimerMessage{Type: timerType, StartTime: startTime, EndTime: endTime, Duration: remaining}
			updateBytes, err := json.Marshal(update)
			if err != nil {
				slog.Error("failed to marshal timer update", "error", err)
				return
			}
			h.BroadcastToGame(gameID, "timer_update", updateBytes)
		}
	}()
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	if gameID == "" {
		gameID = r.URL.Query().Get("game_code")
	}
	if gameID == "" {
		http.Error(w, "game_id or game_code is required", http.StatusBadRequest)
		return
	}

	conn, err := coderws.Accept(w, r, &coderws.AcceptOptions{
		InsecureSkipVerify: true,
		CompressionMode:    coderws.CompressionContextTakeover,
	})
	if err != nil {
		slog.Error("failed to accept websocket", "error", err)
		return
	}
	conn.SetReadLimit(maxMessageSize)

	client := &Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256), GameID: gameID}
	h.hub.Register(client)

	go client.WritePump(r.Context())
	go client.ReadPump(r.Context())
}

func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		c.Hub.unregister <- c
		_ = c.Conn.Close(coderws.StatusNormalClosure, "")
	}()

	for {
		_, data, err := c.Conn.Read(ctx)
		if err != nil {
			return
		}
		if !json.Valid(data) {
			continue
		}
		c.Hub.BroadcastToGame(c.GameID, "client_message", data)
	}
}

func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close(coderws.StatusNormalClosure, "")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-c.Send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.Conn.Write(writeCtx, coderws.MessageText, message)
			cancel()
			if err != nil {
				return
			}
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.Conn.Ping(pingCtx)
			cancel()
			if err != nil {
				return
			}
		}
	}
}
