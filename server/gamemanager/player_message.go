package gamemanager

import "github.com/JStanislav/quoridor-clone/websocket/messages"

type PlayerMessage struct {
	Message messages.IMessage[messages.IncomingMessage]
	IO      *IO
}
