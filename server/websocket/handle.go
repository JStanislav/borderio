package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JStanislav/quoridor-clone/game"
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
	CreateHash func(h string, gameState *game.GameState) *game.GameState
	GetGame    func(h string) *game.GameState
}

func NewHandler(createHash func(h string, gameState *game.GameState) *game.GameState, getGame func(h string) *game.GameState) Handler {
	return Handler{CreateHash: createHash, GetGame: getGame}
}

func (h Handler) Handler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	ppid := r.URL.Query().Get("ppid")
	id := r.PathValue("id")

	var gameState *game.TwoPlayerMatch

	if action == "create" {
		gameState = game.NewTwoPlayerMatch()
		gameState.GameState = *h.CreateHash(id, &gameState.GameState)
		playerOne := player.New(1, ppid, "Player 1", utils.GridPosition{}, 8, utils.Line{}, utils.Line{})
		err := gameState.AddPlayer(playerOne)
		if err != nil {
			fmt.Printf("[ERROR] error adding player to game state, %s\n", err)
			return
		}
	}

	if action == "join" {
		gameState.GameState = *h.GetGame(id)
		playerTwo := player.New(2, ppid, "Player 2", utils.GridPosition{}, 8, utils.Line{}, utils.Line{})
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
	defer c.Close()

	movementsChannel := make(chan player.Play)
	go func() {
		for move := range movementsChannel {
			fmt.Printf("Received move from game state: %+v\n", move)
		}
	}()

	gameState.StartMatch(movementsChannel)

	p1 := gameState.Players[0]
	p2 := gameState.Players[1]

	sendGameState(c, &gameState.GameState, p1, p2)

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
		switch o.PlayerId {
		case 1:
			p = p1
		case 2:
			p = p2
		default:
			fmt.Printf("[ERROR] invalid player ID: %d\n", o.PlayerId)
			continue
		}

		switch o.Type {
		case "playerMove":
			fmt.Printf("Player %d wants to move to row %d, col %d\n", o.PlayerId, o.Target.Row, o.Target.Col)

			err = p.OnPlayerPlay(player.PlayerID(p.ID), player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Row: o.Target.Row, Column: o.Target.Col}})
			if err != nil {
				sendErrorMessage(c, err.Error())
				fmt.Printf("[ERROR] error processing player move, %s\n", err)
				continue
			}
		case "wallPlacement":
			fmt.Printf("Player %d wants to place a wall between [R%d-C%d] and [R%d-C%d] with orientation %s\n", o.PlayerId, o.WallTarget.CellA.Row, o.WallTarget.CellA.Col, o.WallTarget.CellB.Row, o.WallTarget.CellB.Col, o.WallTarget.Orientation)

			err = p.OnPlayerPlay(player.PlayerID(p.ID), player.Play{PlayType: player.WallPlacement, WallPlaced: &utils.WallPosition{CellA: utils.GridPosition{Row: o.WallTarget.CellA.Row, Column: o.WallTarget.CellA.Col}, CellB: utils.GridPosition{Row: o.WallTarget.CellB.Row, Column: o.WallTarget.CellB.Col}}})
			if err != nil {
				sendErrorMessage(c, err.Error())
				fmt.Printf("[ERROR] error processing wall placement, %s\n", err)
				continue
			}
		case "playerReady":
			fmt.Printf("Player %d is ready\n", o.PlayerId)
			p.Ready()
			if gameState.AllPlayersReady() {
				fmt.Println("All players are ready, starting the match")
			}
		}
		// err = c.WriteMessage(mt, []byte("pong"))
		sendGameState(c, &gameState.GameState, p1, p2)
	}
}

func sendGameState(c *websocket.Conn, gameState *game.GameState, p1, p2 *player.Player) {
	gameStateMessage := messages.GameStateStateMessage{
		Type:                "gameState",
		CurrentTurnPlayerId: int(gameState.GetCurrentTurnPlayer().ID),
		PlayerOne: messages.PlayerMessage{
			ID:             int(p1.ID),
			Name:           p1.Name,
			Position:       messages.PositionMessage{Row: p1.Position.Row, Col: p1.Position.Column},
			WallsRemaining: p1.WallsRemaining,
		},
		PlayerTwo: messages.PlayerMessage{
			ID:             int(p2.ID),
			Name:           p2.Name,
			Position:       messages.PositionMessage{Row: p2.Position.Row, Col: p2.Position.Column},
			WallsRemaining: p2.WallsRemaining,
		},
		Walls: gameState.Board.GetWalls(),
	}

	if err := c.WriteJSON(gameStateMessage); err != nil {
		fmt.Printf("[ERROR] error sending game state, %s\n", err)
	}
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
