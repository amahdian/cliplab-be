package svc

import (
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type WebSocketSvc interface {
	Register(userID uuid.UUID, conn *websocket.Conn)
	Send(userID uuid.UUID, event string, data interface{}) error
	Unregister(userID uuid.UUID)
}

var connections map[uuid.UUID]*websocket.Conn

type webSocketSvc struct {
	mu sync.RWMutex
}

func NewWebSocketSvc() WebSocketSvc {
	return &webSocketSvc{}
}

func (h *webSocketSvc) Register(userID uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if connections == nil {
		connections = make(map[uuid.UUID]*websocket.Conn)
	}
	connections[userID] = conn
}

func (h *webSocketSvc) Unregister(userID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(connections, userID)
}

func (h *webSocketSvc) Send(userID uuid.UUID, event string, data interface{}) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conn, ok := connections[userID]
	if !ok {
		return fmt.Errorf("no active connection for user %s", userID)
	}

	msg := map[string]interface{}{
		"event": event,
		"data":  data,
	}
	return conn.WriteJSON(msg)
}
