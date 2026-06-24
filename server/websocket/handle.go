package websocket

import (
	"context"
	"encoding/json"
	"fmt"
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
	id := r.PathValue("id")

	fmt.Print("Received request with action: ", action, " and ppid: ", ppid, " and id: ", id, "\n")

	var gm *gamemanager.GameManager
	var gameState game.TwoPlayerMatch
	var runGame bool

	if action == "create" {
		gameState = *game.NewTwoPlayerMatch()

		timeoutAfterGameOver := h.Context.Value("TimeoutAfterGameOver").(time.Duration)

		gm = gamemanager.NewGameManager(&gameState.GameState, h.UpdateStatsService.UpdateStats, timeoutAfterGameOver)

		err := h.GamesManager.AddGame(id, gm)
		if err != nil {
			fmt.Printf("[ERROR] error creating hash, %s\n", err)

			return

		}

		runGame = true
	}

	if action == "join" {
		gm = h.GamesManager.GetGame(id)
		if gm == nil {
			fmt.Printf("[ERROR] game not found\n")
			return
		}

		gs := gm.Game

		gameState.GameState = *gs
	}

	name := fmt.Sprintf("[PPID: %s]", ppid)
	p := player.New(ppid, name, utils.GridPosition{}, 8, utils.Line{}, utils.Line{})
	err := gameState.AddPlayer(p)
	if err != nil {
		fmt.Printf("[ERROR] error adding player to game state, %s\n", err)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[ERROR] error upgrading, %s\n", err)
		return
	}

	io := gamemanager.NewIO(ppid, c)
	gm.AddPlayer(io)

	if runGame {
		go gm.Run()
	}
}

func (h Handler) GamePing(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

	hash := r.PathValue("hash")

	fmt.Printf("Pinged for game %s\n", hash)

	if h.GamesManager.GetGame(hash) == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h Handler) GamesList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

	games := make([]GameDTO, 0)
	for hash, gm := range h.GamesManager.GetGamesList() {
		games = append(games, GetGameDTO(hash, gm))
	}

	if err := json.NewEncoder(w).Encode(games); err != nil {
		fmt.Printf("[ERROR] error encoding games list, %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
