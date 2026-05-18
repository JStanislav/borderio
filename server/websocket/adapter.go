package websocket

import (
	"github.com/JStanislav/quoridor-clone/gamemanager"
	"github.com/gorilla/websocket"
)

func GetConnectionAdapter(c *websocket.Conn) gamemanager.IOConnection {
	return IOConnection{Conn: c}
}

type IOConnection struct {
	Conn *websocket.Conn
}

func (conn IOConnection) Close() error {
	return conn.Conn.Close()
}

func (conn IOConnection) SendJSON(message any) error {
	return conn.Conn.WriteJSON(message)
}

func (conn IOConnection) SendRaw(message string) error {
	return conn.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}
