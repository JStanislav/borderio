package gamemanager

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/JStanislav/quoridor-clone/websocket/messages"
	"github.com/gorilla/websocket"
)

const sendBufferSize = 32

type IO struct {
	ID   string
	conn *websocket.Conn

	send chan messages.OMessage
}

func NewIO(id string, conn *websocket.Conn) *IO {
	return &IO{
		ID:   id,
		conn: conn,
		send: make(chan messages.OMessage, sendBufferSize),
	}
}

func (io *IO) Send(msg messages.OMessage) {
	select {
	case io.send <- msg:
	default:
		slog.Warn("send buffer full, discarding message", "player_id", io.ID, "message_type", msg.Type)
	}
}

func (io *IO) readPump(inbound chan<- PlayerMessage, done chan<- *IO) {
	defer func() {
		done <- io
	}()

	for {
		_, message, err := io.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			slog.Error("[io] error reading message", "player_id", io.ID, "error", err)
			break
		}

		var o messages.IMessage[messages.IncomingMessage]
		if err = json.Unmarshal(message, &o); err != nil {
			slog.Error("[io] error unmarshaling message", "player_id", io.ID, "error", err)
			break
		}
		inbound <- PlayerMessage{Message: o, IO: io}
	}
}

func (io *IO) writePump() {
	defer func() {
		slog.Info("Closing connection for ppid", "player_id", io.ID)
		now := time.Now().Add(time.Second * 1)
		io.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), now)
	}()

	for msg := range io.send {
		if err := io.conn.WriteJSON(msg); err != nil {
			slog.Error("[io] error writing message", "player_id", io.ID, "error", err)
			break
		}
	}
}
