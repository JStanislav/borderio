package gamemanager

import (
	"errors"

	"github.com/JStanislav/quoridor-clone/game"
	"github.com/gorilla/websocket"
)

type GameManager struct {
	Game *game.GameState

	// Every websocket connection, where key is private player id
	Connections map[string]*websocket.Conn
}

func NewGameManager(game *game.GameState) *GameManager {
	return &GameManager{
		Game:        game,
		Connections: make(map[string]*websocket.Conn),
	}
}

func (gm *GameManager) AddConnection(ppid string, conn *websocket.Conn) {
	gm.Connections[ppid] = conn
}

func (gm *GameManager) RemoveConnection(ppid string) {
	delete(gm.Connections, ppid)
}

func (gm *GameManager) GetConnection(ppid string) *websocket.Conn {
	return gm.Connections[ppid]
}

func (gm *GameManager) Broadcast(msg string) {
	for _, conn := range gm.Connections {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func (gm *GameManager) BroadcastJSON(msg any) {
	for _, conn := range gm.Connections {
		conn.WriteJSON(msg)
	}
}

type Games map[string]*GameManager

func NewGames() Games {
	return make(map[string]*GameManager)
}

func (g *Games) AddGame(h string, gm *GameManager) error {
	if _, ok := (*g)[h]; !ok {
		(*g)[h] = gm
		return nil
	}

	return errors.New("game already exists")
}

func (g Games) GetGame(h string) *GameManager {
	return g[h]
}
