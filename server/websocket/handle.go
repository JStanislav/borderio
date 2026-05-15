package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	GamesManager *gamemanager.Games
}

func NewHandler(gamesManager *gamemanager.Games) Handler {
	return Handler{
		GamesManager: gamesManager,
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
		err := h.GamesManager.AddGame(id, gamemanager.NewGameManager(&gameState.GameState))
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
		gs := h.GamesManager.GetGame(id).Game
		if gs == nil {
			fmt.Printf("[ERROR] game not found\n")
			return
		}
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
	defer c.Close()

	h.GamesManager.GetGame(id).AddConnection(ppid, c)

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

	sendPlayerConfiguration(c, currentPlayer)
	sendLobbyMessage(c, gameState.GameState.Players)

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
			fmt.Printf("Player %d toggled readiness\n", o.PlayerId)
			p.ToggleReady()
			if gameState.AllPlayersReady() && gameState.GameState.PlayerCount == len(*gameState.GameState.Players) {
				fmt.Println("All players are ready, starting the match")
			}
			// sendLobbyMessage(c, gameState.GameState.Players)

			h.GamesManager.GetGame(id).BroadcastJSON(getLobbyMessage(gameState.GameState.Players))
			continue
		}

		// // err = c.WriteMessage(mt, []byte("pong"))
		// sendGameState(c, &gameState.GameState, p1, p2)
	}
}

func sendPlayerConfiguration(c *websocket.Conn, player *player.Player) {
	playerMessage := messages.PlayerConfigurationMessage{
		Type:            "playerConfiguration",
		ID:              int(player.ID),
		Name:            player.Name,
		PrivatePlayerId: player.PrivatePlayerID,
	}

	if err := c.WriteJSON(playerMessage); err != nil {
		fmt.Printf("[ERROR] error sending player configuration, %s\n", err)
	}
}

func getLobbyMessage(players *[]*player.Player) messages.LobbyMessage {
	playersMsg := make([]messages.PlayerMessage, len(*players))
	for i, p := range *players {
		playersMsg[i] = messages.PlayerMessage{
			ID:    int(p.ID),
			Name:  p.Name,
			Ready: p.Ready,
		}
	}
	lobbyMessage := messages.LobbyMessage{
		Type:    "lobby",
		Players: playersMsg,
	}

	return lobbyMessage
}

func sendLobbyMessage(c *websocket.Conn, players *[]*player.Player) {
	lobbyMessage := getLobbyMessage(players)
	if err := c.WriteJSON(lobbyMessage); err != nil {
		fmt.Printf("[ERROR] error sending lobby message, %s\n", err)
	}
}

func sendGameState(c *websocket.Conn, gameState *game.GameState, p1, p2 *player.Player) {
	var currentTurn int
	var walls []utils.WallPosition

	// Match started
	if gameState.StartTime != nil {
		currentTurn = int(gameState.GetCurrentTurnPlayer().ID)
		walls = gameState.Board.GetWalls()
	}

	gameStateMessage := messages.GameStateStateMessage{
		Type:                "gameState",
		CurrentTurnPlayerId: currentTurn,
		PlayerOne: messages.PlayerMessage{
			ID:             int(p1.ID),
			Name:           p1.Name,
			Position:       messages.PositionMessage{Row: p1.Position.Row, Col: p1.Position.Column},
			WallsRemaining: p1.WallsRemaining,
			Ready:          p1.Ready,
		},
		PlayerTwo: messages.PlayerMessage{
			ID:             int(p2.ID),
			Name:           p2.Name,
			Position:       messages.PositionMessage{Row: p2.Position.Row, Col: p2.Position.Column},
			WallsRemaining: p2.WallsRemaining,
			Ready:          p2.Ready,
		},
		Walls: walls,
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
