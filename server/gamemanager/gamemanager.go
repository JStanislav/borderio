package gamemanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/JStanislav/quoridor-clone/external"
	"github.com/JStanislav/quoridor-clone/game"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/websocket/messages"
)

type GameManager struct {
	Game *game.GameState

	gameOver bool

	// Every websocket connection, where key is private player id
	IOManager IOManager

	UpdateStats external.UpdateStats
}

func NewGameManager(game *game.GameState, ioManager IOManager, updateStats external.UpdateStats) *GameManager {
	return &GameManager{
		Game:        game,
		gameOver:    false,
		IOManager:   ioManager,
		UpdateStats: updateStats,
	}
}

func (gm *GameManager) AddConnection(ppid string, conn IOConnection) {
	gm.IOManager.AddConnection(ppid, conn)

}

func (gm *GameManager) RemoveConnection(ppid string) {
	gm.IOManager.RemoveConnection(ppid)
}

func (gm *GameManager) Broadcast(msg string) error {
	err := gm.IOManager.Broadcast(msg)
	if err != nil {
		fmt.Println("error broadcasting ", err)
		errMessage := fmt.Sprintf("error sending %s", msg)
		return errors.New(errMessage)
	}
	return nil
}

func (gm *GameManager) BroadcastJSON(msg any) error {
	err := gm.IOManager.BroadcastJSON(msg)
	if err != nil {
		fmt.Println("error broadcasting JSON", err)
		errMessage := fmt.Sprintf("error sending %s", msg)
		return errors.New(errMessage)
	}
	return nil
}

// Function to broadcast messages to every player except the specified in the parameter
func (gm *GameManager) BroadcastExcept(msg any, ppids []string) error {
	err := gm.IOManager.BroadcastJSONExcept(msg, ppids)
	if err != nil {
		fmt.Println("error broadcasting JSON Except", err)
		errMessage := fmt.Sprintf("error sending %s", msg)
		return errors.New(errMessage)
	}
	return nil
}

func (gm *GameManager) BroadcastGameState() error {
	gameStateMesage := messages.GetGameStateMessage(gm.Game)
	err := gm.BroadcastJSON(gameStateMesage)
	if err != nil {
		return err
	}
	return nil
}

func (gm *GameManager) CleanUpConnection(ppid string) {
	c := gm.IOManager.GetConnection(ppid)
	c.Close()
	gm.RemoveConnection(ppid)
}

func (gm *GameManager) PlayerLeft(player player.Player) {
	gm.Game.RemovePlayer(player.ID)
	gm.CleanUpConnection(player.PrivatePlayerID)
	gm.BroadcastJSON(messages.GetPlayerLeftMessage(player))

	// Broadcast lobby, what would happen if player left in middle of the game?
	gm.BroadcastJSON(messages.GetLobbyMessage(gm.Game.Players))
}

func (gm *GameManager) PlayerJoined(player player.Player) {
	message := messages.GetJoinedMessage(player)
	gm.BroadcastExcept(message, []string{player.PrivatePlayerID})

	gm.SyncLobbyState()
	gm.SyncPlayerConfiguration(player)
	gm.SyncMatchConfiguration(player)
}

func (gm *GameManager) SyncMatchConfiguration(player player.Player) {
	playerConn := gm.IOManager.GetConnection(player.PrivatePlayerID)
	message := messages.GetMatchConfigurationMessage(gm.Game.PlayerCount)
	if err := playerConn.SendJSON(message); err != nil {
		fmt.Printf("[ERROR] error sending match configuration, %s\n", err)
	}
}

func (gm *GameManager) GameOver() {
	gm.gameOver = true
	gm.BroadcastJSON(messages.GetLobbyMessage(gm.Game.Players))

	err := gm.UpdateStats(gm.Game.GetGameStats())
	if err != nil {
		fmt.Printf("[ERROR] error updating stats, %s\n", err)
	}

	time.AfterFunc(3*time.Minute, func() {
		fmt.Printf("closing all connections\n")
		gm.DisconnectAll()
	})
}

func (gm *GameManager) DisconnectAll() {
	err := gm.IOManager.DisconnectAll()
	if err != nil {
		fmt.Printf("[ERROR] error disconnecting all connections, %s\n", err)
	}
}

func (gm *GameManager) IsGameOver() bool {
	return gm.gameOver
}

func (gm *GameManager) SyncLobbyState() {
	gm.BroadcastJSON(messages.GetLobbyMessage(gm.Game.Players))
}

func (gm *GameManager) SyncPlayerConfiguration(player player.Player) {
	playerConn := gm.IOManager.GetConnection(player.PrivatePlayerID)

	message := messages.PlayerConfigurationMessage{
		Type:            "playerConfiguration",
		ID:              int(player.ID),
		Name:            player.Name,
		PrivatePlayerId: player.PrivatePlayerID,
	}

	if err := playerConn.SendJSON(message); err != nil {
		fmt.Printf("[ERROR] error sending player configuration, %s\n", err)
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

func (g *Games) RemoveGame(h string) {
	delete(*g, h)
}
