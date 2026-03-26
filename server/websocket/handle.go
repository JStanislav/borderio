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

func Handle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[ERROR] error upgrading, %s\n", err)
		return
	}
	defer c.Close()

	p1 := &player.Player{ID: 1, Name: "Colo"}
	p2 := &player.Player{ID: 2, Name: "Stan"}

	gameState := game.New()

	movementsChannel := make(chan player.Play)
	go func() {
		for move := range movementsChannel {
			fmt.Printf("Received move from game state: %+v\n", move)
		}
	}()

	gameState.StartMatch(p1, p2, movementsChannel)
	sendGameState(c, gameState, p1, p2)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("[ERROR] error reading message, %s\n", err)
			break
		}

		var o messages.Message
		if err = json.Unmarshal(message, &o); err != nil {
			fmt.Printf("[ERROR] error unmarshaling message, %s\n", err)
			break
		}

		switch o.Type {
		case "playerMove":
			fmt.Printf("Player %d wants to move to row %d, col %d\n", o.PlayerId, o.Target.Row, o.Target.Col)
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
			err = p.OnPlayerPlay(player.PlayerID(p.ID), player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Row: o.Target.Row, Column: o.Target.Col}})
			if err != nil {
				sendErrorMessage(c, err.Error())
				fmt.Printf("[ERROR] error processing player move, %s\n", err)
				continue
			}
		case "wallPlacement":
			fmt.Printf("Player %d wants to place a wall between [R%d-C%d] and [R%d-C%d] with orientation %s\n", o.PlayerId, o.WallTarget.CellA.Row, o.WallTarget.CellA.Col, o.WallTarget.CellB.Row, o.WallTarget.CellB.Col, o.WallTarget.Orientation)
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
			err = p.OnPlayerPlay(player.PlayerID(p.ID), player.Play{PlayType: player.WallPlacement, WallPlaced: &utils.WallPosition{CellA: utils.GridPosition{Row: o.WallTarget.CellA.Row, Column: o.WallTarget.CellA.Col}, CellB: utils.GridPosition{Row: o.WallTarget.CellB.Row, Column: o.WallTarget.CellB.Col}}})
			if err != nil {
				sendErrorMessage(c, err.Error())
				fmt.Printf("[ERROR] error processing wall placement, %s\n", err)
				continue
			}
		}
		// err = c.WriteMessage(mt, []byte("pong"))
		sendGameState(c, gameState, p1, p2)
	}
}

func sendGameState(c *websocket.Conn, gameState *game.GameState, p1, p2 *player.Player) {
	gameStateMessage := messages.GameStateStateMessage{
		Type: "gameState",
		PlayerOne: messages.PlayerMessage{
			ID:       int(p1.ID),
			Name:     p1.Name,
			Position: messages.PositionMessage{Row: p1.Position.Row, Col: p1.Position.Column},
		},
		PlayerTwo: messages.PlayerMessage{
			ID:       int(p2.ID),
			Name:     p2.Name,
			Position: messages.PositionMessage{Row: p2.Position.Row, Col: p2.Position.Column},
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
