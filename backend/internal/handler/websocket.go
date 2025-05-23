package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	ws "github.com/zizouhuweidi/dahaa/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub *ws.Hub
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleWebSocket handles incoming WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	// Get game ID from query parameter
	gameID := c.QueryParam("game_id")
	if gameID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "game_id is required",
		})
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}

	// Create new client
	client := &ws.Client{
		Hub:    h.hub,
		Conn:   conn,
		GameID: gameID,
		Send:   make(chan []byte, 256),
	}

	// Register client
	h.hub.Register(client)

	// Start goroutines for reading and writing
	go client.ReadPump()
	go client.WritePump()

	return nil
}
