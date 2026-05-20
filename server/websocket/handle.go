package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JStanislav/quoridor-clone/external"
	"github.com/JStanislav/quoridor-clone/game"
	"github.com/JStanislav/quoridor-clone/gamemanager"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
	"github.com/JStanislav/quoridor-clone/websocket/messages"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
}

type Handler struct {
	GamesManager       *gamemanager.Games
	UpdateStatsService external.UpdateStatsService
}

func NewHandler(gamesManager *gamemanager.Games, updateStatsService external.UpdateStatsService) Handler {
	return Handler{
		GamesManager:       gamesManager,
		UpdateStatsService: updateStatsService,
	}
}

func (h Handler) Handler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	ppid := r.URL.Query().Get("ppid")
	id := r.PathValue("id")

	fmt.Print("Received request with action: ", action, " and ppid: ", ppid, " and id: ", id, "\n")

	gameState := game.NewTwoPlayerMatch()
	var currentPlayer *player.Player

	if action == "create" {
		cm := gamemanager.NewConnectionsManager()
		err := h.GamesManager.AddGame(id, gamemanager.NewGameManager(&gameState.GameState, cm, h.UpdateStatsService.UpdateStats))
		if err != nil {
			fmt.Printf("[ERROR] error creating hash, %s\n", err)
			return
		}

		playerOne := player.New(1, ppid, "Player 1", utils.GridPosition{}, 8, utils.Line{}, utils.Line{})
		currentPlayer = playerOne
		err = gameState.AddPlayer(playerOne)
		if err != nil {
			fmt.Printf("[ERROR] error adding player to game state, %s\n", err)
			return
		}
	}

	if action == "join" {
		gm := h.GamesManager.GetGame(id)
		if gm == nil {
			fmt.Printf("[ERROR] game not found\n")
			return
		}
		gs := gm.Game

		gameState.GameState = *gs

		playerTwo := player.New(2, ppid, "Player 2", utils.GridPosition{}, 8, utils.Line{}, utils.Line{})
		currentPlayer = playerTwo
		err := gameState.AddPlayer(playerTwo)
		if err != nil {
			fmt.Printf("[ERROR] error adding player to game state, %s\n", err)
			return
		}
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[ERROR] error upgrading, %s\n", err)
		return
	}
	defer func() {
		h.GamesManager.GetGame(id).PlayerLeft(*currentPlayer)
	}()

	ioConn := GetConnectionAdapter(c)

	h.GamesManager.GetGame(id).AddConnection(ppid, ioConn)

	// This movements channel has to go away from here. It should be only one in the game, not one for every connection
	movementsChannel := make(chan player.Play)
	go func() {
		for move := range movementsChannel {
			fmt.Printf("Received move from game state: %+v\n", move)
		}
	}()

	lobbyChannel := make(chan messages.LobbyMessage)
	go func() {
		for lobbyMessage := range lobbyChannel {
			fmt.Printf("Received lobby message%+v\n", lobbyMessage)
		}
	}()

	h.GamesManager.GetGame(id).PlayerJoined(*currentPlayer)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("[ERROR] error reading message, %s\n", err)
			break
		}

		var o messages.Message[any]
		if err = json.Unmarshal(message, &o); err != nil {
			fmt.Printf("[ERROR] error unmarshaling message, %s\n", err)
			break
		}

		var p *player.Player
		p = gameState.GetPlayerPPID(o.PrivatePlayerId)
		if p == nil {
			fmt.Printf("[ERROR] player with ppid %s not found\n", o.PrivatePlayerId)
			continue
		}

		switch o.Type {
		case "startGame":
			fmt.Printf("Player %d wants to start the game\n", p.ID)

			gameState.StartMatch(movementsChannel)

			h.GamesManager.GetGame(id).BroadcastGameState()
		case "playerMove":
			fmt.Printf("Player %d wants to move to row %d, col %d\n", p.ID, o.Target.Row, o.Target.Col)

			err = p.OnPlayerPlay(player.PlayerID(p.ID), player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Row: o.Target.Row, Column: o.Target.Col}})
			if err != nil {
				sendErrorMessage(c, err.Error())
				fmt.Printf("[ERROR] error processing player move, %s\n", err)
				continue
			}

			h.GamesManager.GetGame(id).BroadcastGameState()

			if p.IsWinner() {
				fmt.Printf("Player %d wins!\n", p.ID)
				h.GamesManager.GetGame(id).GameOver()
			}
		case "wallPlacement":
			fmt.Printf("Player %d wants to place a wall between [R%d-C%d] and [R%d-C%d] with orientation %s\n", p.ID, o.WallTarget.CellA.Row, o.WallTarget.CellA.Col, o.WallTarget.CellB.Row, o.WallTarget.CellB.Col, o.WallTarget.Orientation)

			err = p.OnPlayerPlay(player.PlayerID(p.ID), player.Play{PlayType: player.WallPlacement, WallPlaced: &utils.WallPosition{CellA: utils.GridPosition{Row: o.WallTarget.CellA.Row, Column: o.WallTarget.CellA.Col}, CellB: utils.GridPosition{Row: o.WallTarget.CellB.Row, Column: o.WallTarget.CellB.Col}}})
			if err != nil {
				sendErrorMessage(c, err.Error())
				fmt.Printf("[ERROR] error processing wall placement, %s\n", err)
				continue
			}

			h.GamesManager.GetGame(id).BroadcastGameState()
		case "playerReady":
			fmt.Printf("Player %d toggled readiness\n", p.ID)
			p.ToggleReady()
			if gameState.AllPlayersReady() && gameState.GameState.PlayerCount == len(*gameState.GameState.Players) {
				fmt.Println("All players are ready, starting the match")
			}

			h.GamesManager.GetGame(id).SyncLobbyState()
			continue
		}
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

func sendErrorMessage(c *websocket.Conn, errorMessage string) {
	errorMessageStruct := messages.ErrorMessage{
		Type:    "error",
		Message: errorMessage,
	}
	if err := c.WriteJSON(errorMessageStruct); err != nil {
		fmt.Printf("[ERROR] error sending error message, %s\n", err)
	}
}
