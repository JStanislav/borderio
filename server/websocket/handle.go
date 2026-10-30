package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/JStanislav/quoridor-clone/external"
	"github.com/JStanislav/quoridor-clone/game"
	"github.com/JStanislav/quoridor-clone/gamemanager"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
}

type Handler struct {
	Context            context.Context
	GamesManager       *gamemanager.Games
	UpdateStatsService external.UpdateStatsService
}

func NewHandler(ctx context.Context, gamesManager *gamemanager.Games, updateStatsService external.UpdateStatsService) Handler {
	return Handler{
		Context:            ctx,
		GamesManager:       gamesManager,
		UpdateStatsService: updateStatsService,
	}
}

func (h Handler) StartOrJoinGame(w http.ResponseWriter, r *http.Request) {
	ppid := r.URL.Query().Get("ppid")
	gameId := r.PathValue("game_id")
	name := r.URL.Query().Get("name")

	nameWithID := fmt.Sprintf("%s [%s]", name, ppid)
	slog.Info("Received request", "game", gameId, "player", nameWithID)

	if err := ValidateParams(ppid, name, gameId); err != nil {
		slog.Error("invalid request", "error", err, "game", gameId, "player", name)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	create := h.GamesManager.GetGame(gameId) == nil

	ctx := context.WithValue(context.Background(), "game", gameId)

	addPlayer := false

	var match game.TwoPlayerMatch
	var gm *gamemanager.GameManager

	if create {
		slog.Info("Creating game", "game", gameId, "player", nameWithID)

		match = *game.NewTwoPlayerMatch()

		timeoutAfterGameOver := h.Context.Value("TimeoutAfterGameOver").(time.Duration)

		gm = gamemanager.NewGameManager(ctx, gameId, &match.GameState, h.UpdateStatsService.UpdateStats, timeoutAfterGameOver)
		if gracePeriod, ok := h.Context.Value("ReconnectionGracePeriod").(time.Duration); ok {
			gm.ReconnectionGracePeriod = gracePeriod
		}

		err := h.GamesManager.AddGame(gameId, gm)
		if err != nil {
			slog.Error("error creating game", "error", err, "game", gameId, "player", nameWithID)
			return
		}

		addPlayer = true
		defer func() {
			go gm.Run()
		}()

	} else {
		slog.Info("Joining game", "game", gameId, "player", nameWithID)

		gm = h.GamesManager.GetGame(gameId)
		if gm == nil {
			slog.Error("unexpected error: game not found", "game", gameId, "player", nameWithID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if gm.Game.GetPlayerByPPID(ppid) != nil {
			// A player already in the game is coming back (e.g. page reload or reconnection during the grace period)
			slog.Info("Player rejoining game", "game", gameId, "player", nameWithID)
		} else if len(*gm.Game.Players) < 2 {
			// Joining the game as a player
			addPlayer = true
		}

		gs := gm.Game

		match.GameState = *gs
	}

	if addPlayer {
		err := match.AddPlayer(name, ppid, 8)
		if err != nil {
			slog.Error("error adding player to game state", "error", err, "game", gameId, "player", nameWithID)
			return
		}
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("error upgrading connection", "error", err, "game", gameId, "player", nameWithID)
		return
	}

	io := gamemanager.NewIO(ppid, name, c)
	gm.AddPlayer(io)

}

func (h Handler) GamePing(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")

	slog.Info("Received ping request", "game", hash)

	if h.GamesManager.GetGame(hash) == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h Handler) GamesList(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received games list request")

	games := make([]GameDTO, 0)
	for hash, gm := range h.GamesManager.GetGamesList() {
		games = append(games, GetGameDTO(hash, gm))
	}

	if err := json.NewEncoder(w).Encode(games); err != nil {
		slog.Error("error encoding games list", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received health check request")

	w.WriteHeader(http.StatusNoContent)
}
