package gamemanager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/JStanislav/quoridor-clone/external"
	"github.com/JStanislav/quoridor-clone/game"
	"github.com/JStanislav/quoridor-clone/gamemanager/lobby"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
	"github.com/JStanislav/quoridor-clone/websocket/messages"
)

type GameManager struct {
	ID   string
	Game *game.GameState

	Started      bool
	gameOver     bool
	GameTimedOut bool

	UpdateStats external.UpdateStats

	TimeoutAfterGameOver time.Duration

	IOs []*IO

	inbound chan PlayerMessage
	join    chan *IO
	leave   chan *IO
	quit    chan struct{}

	Lobbyists map[string]*lobby.Lobbyist
}

func NewGameManager(context context.Context, id string, game *game.GameState, updateStats external.UpdateStats, timeoutAfterGameOver time.Duration) *GameManager {
	return &GameManager{
		ID:                   id,
		Game:                 game,
		gameOver:             false,
		UpdateStats:          updateStats,
		GameTimedOut:         false,
		TimeoutAfterGameOver: timeoutAfterGameOver,

		inbound: make(chan PlayerMessage),
		join:    make(chan *IO, 8),
		leave:   make(chan *IO, 8),
		quit:    make(chan struct{}),

		Lobbyists: make(map[string]*lobby.Lobbyist),
	}
}

func (gm *GameManager) AddPlayer(p *IO) {
	gm.join <- p
	go p.readPump(gm.inbound, gm.leave)
	go p.writePump()
}

func (gm *GameManager) Stop() {
	gm.broadcastJSON(messages.OMessage{Type: "gameTimedOut"})
	gm.GameTimedOut = true
	close(gm.quit)
}

func (gm *GameManager) Run() {
	slog.Info("game manager started", "game_id", gm.ID)
	defer slog.Info("game manager stopped", "game_id", gm.ID)

	for {
		select {
		case <-gm.quit:
			gm.shutDownAll()
			return

		case p := <-gm.join:
			gm.handleJoin(p)

		case p := <-gm.leave:
			gm.handleLeave(p)

		case msg := <-gm.inbound:
			gm.handleMessage(msg)
		}
	}
}

func (gm *GameManager) shutDownAll() {
	for _, p := range gm.IOs {
		close(p.send)
	}
}

func (gm *GameManager) handleJoin(io *IO) {
	// if gm.Game.StartTime != nil {
	// 	io.Send(messages.GetAlreadyStartedMessage())
	// 	close(io.send)
	// 	return
	// }

	gm.IOs = append(gm.IOs, io)
	slog.Info("[manager] player joined", "game_id", gm.ID, "player_count", len(gm.IOs), "player_id", io.ID)

	var l *lobby.Lobbyist
	p := gm.Game.GetPlayerByPPID(io.ID)

	if p == nil {
		// player is surely a spectator
		l = lobby.NewLobbyist(-1, io.Name)
		l.Spectator = true
	} else {
		l = lobby.NewLobbyist(int(p.ID), p.Name)
		if len(gm.Lobbyists) == 0 {
			l.Host = true
		}
	}
	l.Connected = true

	gm.Lobbyists[io.ID] = l

	message := messages.GetJoinedMessage(*l)
	gm.broadcastExcept(message, []string{io.ID})

	gm.syncLobbyState()
	gm.syncPlayerConfiguration(io)
	gm.syncMatchConfiguration(io)
	if gm.Started {
		gm.broadcastGameState()
	}
}

func (gm *GameManager) handleLeave(io *IO) {
	idx := gm.indexOf(io)
	if idx == -1 {
		return
	}

	slog.Info("[manager] player left", "game_id", gm.ID, "player_id", io.ID)
	gm.IOs = append(gm.IOs[:idx], gm.IOs[idx+1:]...)
	close(io.send)

	l := *gm.Lobbyists[io.ID]
	l.Connected = false

	if !gm.Started && !l.Spectator {
		// players are only removed from the game when the game is not started
		// so the game can remember players positions one the game has started
		p := gm.Game.GetPlayerByPPID(io.ID)
		gm.Game.RemovePlayer(p.ID)
	}

	if l.Spectator || !gm.Started {
		gm.RemoveLobbyist(io.ID)
	}

	gm.broadcastJSON(messages.GetPlayerLeftMessage(l))
	gm.syncLobbyState()

	// TODO: implement when player abandons in the middle of the game
	// if gm.Game.StartTime != nil && !gm.gameOver {
	// 	if len(gm.IOs) < gm.Game.PlayerCount {
	// 		gm.endGame(fmt.Sprintf("player %s abandoned", p.Name))
	// 		return
	// 	}
	// 	gm.broadcastGameState()
	// }
}

func (gm *GameManager) RemoveLobbyist(ppid string) {
	wasHost := gm.Lobbyists[ppid].Host
	delete(gm.Lobbyists, ppid)
	if wasHost && len(gm.Lobbyists) > 0 {
		for _, l := range gm.Lobbyists {
			l.Host = true
			break
		}
	}
}

func (gm *GameManager) handleMessage(msg PlayerMessage) {
	var p *player.Player
	p = gm.Game.GetPlayerByPPID(msg.Message.PrivatePlayerId)
	name := fmt.Sprintf("%s [%s]", p.Name, p.PrivatePlayerID)

	if p == nil {
		slog.Error("[manager] player not found", "game_id", gm.ID, "player", name)
		return
	}

	switch msg.Message.Type {
	case "startGame":
		if gm.Started {
			slog.Error("[manager] game is already started", "game_id", gm.ID, "player", name)
			msg.IO.Send(messages.OMessage{
				Type:    "error",
				Payload: "game is already started",
			})
			return
		}

		if !gm.Lobbyists[msg.IO.ID].Host {
			slog.Error("[manager] only the host can start the game", "game_id", gm.ID, "player", name)
			msg.IO.Send(messages.OMessage{
				Type:    "error",
				Payload: "only the host can start the game",
			})
			return
		}

		if !gm.AllPlayersReady() {
			slog.Error("[manager] not all players are ready", "game_id", gm.ID, "player", name)
			msg.IO.Send(messages.OMessage{
				Type:    "error",
				Payload: "not all players are ready",
			})
			return
		}

		gm.Game.StartMatchWithMovementsChannel()

		gm.Started = true

		gm.syncLobbyState()
		gm.broadcastGameState()
	case "playerMove":
		mov := msg.Message.Payload

		if gm.IsGameOver() {
			slog.Error("[manager] game is already over, cannot make a move", "game_id", gm.ID, "player", name)
			msg.IO.Send(messages.OMessage{
				Type:    "error",
				Payload: "game is already over",
			})
			return
		}

		err := p.OnPlayerPlay(player.PlayerID(p.ID), player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Row: mov.Target.Row, Column: mov.Target.Col}})
		if err != nil {
			msg.IO.Send(messages.OMessage{
				Type:    "error",
				Payload: err.Error(),
			})
			slog.Warn("[manager] error processing player move", "game_id", gm.ID, "player", name, "error", err)
			return
		}

		gm.broadcastGameState()

		if p.IsWinner() {
			slog.Info("[manager] player wins", "game_id", gm.ID, "player", name)
			gm.endGame(fmt.Sprintf("player %s wins", p.Name))
		}
	case "wallPlacement":
		wallMsg := msg.Message.Payload

		if gm.IsGameOver() {
			slog.Error("[manager] game is already over, cannot make a move", "game_id", gm.ID, "player", name)
			msg.IO.Send(messages.OMessage{
				Type:    "error",
				Payload: "game is already over, cannot make a move",
			})
			return
		}

		err := p.OnPlayerPlay(player.PlayerID(p.ID), player.Play{PlayType: player.WallPlacement, WallPlaced: &utils.WallPosition{CellA: utils.GridPosition{Row: wallMsg.WallTarget.CellA.Row, Column: wallMsg.WallTarget.CellA.Col}, CellB: utils.GridPosition{Row: wallMsg.WallTarget.CellB.Row, Column: wallMsg.WallTarget.CellB.Col}}})
		if err != nil {
			msg.IO.Send(messages.OMessage{
				Type:    "error",
				Payload: err.Error(),
			})
			slog.Warn("[manager] error processing wall placement", "game_id", gm.ID, "player", name, "error", err)
			return
		}

		gm.broadcastGameState()
	case "playerReady":
		slog.Info("[manager] player readiness toggled", "game_id", gm.ID, "player", name)
		gm.Lobbyists[p.PrivatePlayerID].ToggleReady()

		gm.syncLobbyState()
		return
	}
}

func (gm *GameManager) broadcastJSON(msg messages.OMessage) error {
	for _, io := range gm.IOs {
		io.Send(msg)
	}

	return nil
}

// Function to broadcast messages to every player except the specified in the parameter
func (gm *GameManager) broadcastExcept(msg messages.OMessage, ppids []string) error {
	for _, io := range gm.IOs {
		if !slices.Contains(ppids, io.ID) {
			io.Send(msg)
		}
	}
	return nil
}

func (gm *GameManager) broadcastGameState() error {
	gameStateMesage := messages.GetGameStateMessage(gm.Game)
	return gm.broadcastJSON(gameStateMesage)
}

func (gm *GameManager) syncMatchConfiguration(io *IO) {
	message := messages.GetMatchConfigurationMessage(gm.Game.PlayersCapacity)
	io.Send(message)
}

func (gm *GameManager) endGame(reason string) {
	gm.gameOver = true
	gm.syncLobbyState()

	err := gm.UpdateStats(gm.Game.GetGameStats())
	if err != nil {
		slog.Error("[manager] error updating stats", "game_id", gm.ID, "error", err)
	}

	time.AfterFunc(gm.TimeoutAfterGameOver, func() {
		slog.Info("[manager] closing all connections", "game_id", gm.ID, "reason", reason)
		gm.Stop()
	})

	notificationTimeSpan := 10 // in seconds

	time.AfterFunc(gm.TimeoutAfterGameOver-time.Duration(notificationTimeSpan)*time.Second, func() {
		gm.broadcastJSON(messages.GetWillTimeOutMessage(notificationTimeSpan))
	})
}

func (gm *GameManager) IsGameOver() bool {
	return gm.gameOver
}

func (gm *GameManager) AllPlayersReady() bool {
	for _, p := range gm.Lobbyists {
		if !p.Ready && !p.Spectator {
			return false
		}
	}
	return true
}

func (gm *GameManager) syncLobbyState() {
	winnerId := gm.Game.GetWinnerPlayerId()
	gm.broadcastJSON(messages.GetLobbyMessage(gm.Started, gm.Lobbyists, gm.ID, winnerId))
}

func (gm *GameManager) syncPlayerConfiguration(io *IO) {
	p := gm.Game.GetPlayerByPPID(io.ID)

	id := -1
	if p != nil {
		id = int(p.ID)
	}

	message := messages.OMessage{
		Type: "playerConfiguration",
		Payload: messages.PlayerConfigurationMessage{
			ID:              id,
			Name:            io.Name,
			PrivatePlayerId: io.ID,
		},
	}

	io.Send(message)
}

func (gm *GameManager) indexOf(target *IO) int {
	for i, p := range gm.IOs {
		if p.ID == target.ID {
			return i
		}
	}
	return -1
}

type GamesContainer struct {
	Games map[string]*GameManager

	GC *EndedGamesCollector
}

func NewGamesContainer(GCThreshold int) *GamesContainer {
	games := NewGames()
	return &GamesContainer{Games: games, GC: NewGC(&games, GCThreshold)}
}

type Games map[string]*GameManager

func NewGames() Games {
	return make(map[string]*GameManager)
}

func (g *Games) GetGamesList() Games {
	return *g
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

func (g *Games) DeleteOldGames() {
	for h, gm := range g.GetGamesList() {
		remove := false
		if gm.IsGameOver() && (gm.GameTimedOut || len(gm.IOs) == 0) {
			slog.Info("[manager] game timed out", "game_id", h)
			remove = true
		}

		if gm.Game.StartTime != nil && len(gm.IOs) == 0 {
			slog.Info("[manager] game didn't start and no players in lobby", "game_id", h)
			remove = true
		}

		// gm is already stopped (gm.Stop()), when timeout is reached. This step is called strictly after that.
		if remove {
			slog.Info("[manager] deleting game", "game_id", h)
			g.RemoveGame(h)
		}
	}
}
