package websocket

import (
	"sync"
	"time"

	"github.com/JStanislav/quoridor-clone/gamemanager"
	"github.com/gorilla/websocket"
)

func GetConnectionAdapter(c *websocket.Conn) gamemanager.IOConnection {
	return &IOConnection{Conn: c}
}

type IOConnection struct {
	m    sync.RWMutex
	Conn *websocket.Conn
}

func (conn *IOConnection) Close() error {
	conn.m.Lock()
	defer conn.m.Unlock()

	deadline := time.Now().Add(1 * time.Second)
	conn.Conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing gracefully"),
		deadline)

	return nil
}

func (conn *IOConnection) SendJSON(message any) error {
	conn.m.Lock()
	defer conn.m.Unlock()

	return conn.Conn.WriteJSON(message)
}

func (conn *IOConnection) SendRaw(message string) error {
	conn.m.Lock()
	defer conn.m.Unlock()

	return conn.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}
