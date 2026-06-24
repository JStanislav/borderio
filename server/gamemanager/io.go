package gamemanager

import (
	"encoding/json"
	"fmt"

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
		fmt.Printf("[player %s] send buffer full, discarding message\n (type=%s)", io.ID, msg.Type)
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
			fmt.Printf("[ERROR] error reading message, %s\n", err)
			break
		}

		var o messages.IMessage[messages.IncomingMessage]
		if err = json.Unmarshal(message, &o); err != nil {
			fmt.Printf("[ERROR] error unmarshaling message, %s\n", err)
			break
		}
		inbound <- PlayerMessage{Message: o, IO: io}
	}
}

func (io *IO) writePump() {
	defer func() {
		fmt.Printf("Closing connection for ppid: %s\n", io.ID)
		io.conn.Close()
	}()

	for msg := range io.send {
		if err := io.conn.WriteJSON(msg); err != nil {
			fmt.Printf("[ERROR] error writing message, %s\n", err)
			break
		}
	}
}
