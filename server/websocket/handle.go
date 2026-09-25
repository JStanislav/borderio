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
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
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

func (h Handler) Handler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	ppid := r.URL.Query().Get("ppid")
	gameHash := r.PathValue("id")
	name := r.URL.Query().Get("name")

	if err := ValidateParams(action, ppid, name, gameHash); err != nil {
		slog.Error("invalid request", "error", err, "game", gameHash, "player", name)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ctx := context.WithValue(context.Background(), "game", gameHash)
	ctx = context.WithValue(ctx, "player_name", fmt.Sprintf("%s [%s]", name, ppid))

	nameWithID := fmt.Sprintf("%s [%s]", name, ppid)
	slog.Info("Received request", "action", action, "game", gameHash, "player", nameWithID)

	var gm *gamemanager.GameManager
	var gameState game.TwoPlayerMatch
	var runGame bool

	if action == "create" {
		gameState = *game.NewTwoPlayerMatch()

		timeoutAfterGameOver := h.Context.Value("TimeoutAfterGameOver").(time.Duration)

		gm = gamemanager.NewGameManager(ctx, gameHash, &gameState.GameState, h.UpdateStatsService.UpdateStats, timeoutAfterGameOver)

		err := h.GamesManager.AddGame(gameHash, gm)
		if err != nil {
			slog.Error("error creating game", "error", err, "game", gameHash, "player", nameWithID)
			return
		}

		runGame = true
	}

	if action == "join" {
		gm = h.GamesManager.GetGame(gameHash)
		if gm == nil {
			slog.Warn("game not found", "game", gameHash, "player", nameWithID)
			return
		}

		gs := gm.Game

		gameState.GameState = *gs
	}

	if action == "spectate" {
		gm = h.GamesManager.GetGame(gameHash)
		if gm == nil {
			slog.Warn("game not found", "game", gameHash, "player", nameWithID)
			w.WriteHeader(http.StatusNotFound)
			return
		}

	}

	p := player.New(ppid, name, utils.GridPosition{}, 8, utils.Line{}, utils.Line{})
	err := gameState.AddPlayer(p)
	if err != nil {
		slog.Error("error adding player to game state", "error", err, "game", gameHash, "player", nameWithID)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("error upgrading connection", "error", err, "game", gameHash, "player", nameWithID)
		return
	}

	io := gamemanager.NewIO(ppid, c)
	gm.AddPlayer(io)

	if runGame {
		go gm.Run()
	}
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
