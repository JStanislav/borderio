package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/JStanislav/quoridor-clone/game"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Message struct {
	Type     string `json:"type"`
	PlayerId int    `json:"playerId"`
	Target   struct {
		Row int `json:"row"`
		Col int `json:"col"`
	} `json:"target"`
	WallTarget WallTargetMessage `json:"wallTarget"`
}

type WallTargetMessage struct {
	CellA       PositionMessage `json:"cellA"`
	CellB       PositionMessage `json:"cellB"`
	Orientation string          `json:"orientation"` // "horizontal" or "vertical"
}

type PositionMessage struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type PlayerMessage struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Position PositionMessage `json:"position"`
}

type GameStateStateMessage struct {
	Type      string               `json:"type"`
	PlayerOne PlayerMessage        `json:"playerOne"`
	PlayerTwo PlayerMessage        `json:"playerTwo"`
	Walls     []utils.WallPosition `json:"walls"`
}

func sendGameState(c *websocket.Conn, gameState *game.GameState, p1, p2 *player.Player) {
	gameStateMessage := GameStateStateMessage{
		Type: "gameState",
		PlayerOne: PlayerMessage{
			ID:       int(p1.ID),
			Name:     p1.Name,
			Position: PositionMessage{Row: p1.Position.Row, Col: p1.Position.Column},
		},
		PlayerTwo: PlayerMessage{
			ID:       int(p2.ID),
			Name:     p2.Name,
			Position: PositionMessage{Row: p2.Position.Row, Col: p2.Position.Column},
		},
	}
	gameStateMessage.Walls = gameState.Board.GetWalls()

	if err := c.WriteJSON(gameStateMessage); err != nil {
		fmt.Printf("[ERROR] error sending game state, %s\n", err)
	}

}

func handle(w http.ResponseWriter, r *http.Request) {
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

		fmt.Printf("received %s\n", message)
		var o Message
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
				fmt.Printf("[ERROR] error processing wall placement, %s\n", err)
				continue
			}
		}
		// err = c.WriteMessage(mt, []byte("pong"))
		sendGameState(c, gameState, p1, p2)
	}
}

func main() {
	fmt.Println("Hello World")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
