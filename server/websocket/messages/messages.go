package messages

import (
	"github.com/JStanislav/quoridor-clone/game"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
)

type OMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type IMessage[T any] struct {
	Type    string `json:"type"`
	Payload T      `json:"payload"`

	PrivatePlayerId string `json:"ppid"`
}

type IncomingMessage struct {
	Target     PositionMessage   `json:"target"`     // remove this and put it in Payload when we have time to refactor the client
	WallTarget WallTargetMessage `json:"wallTarget"` // remove this and put it in Payload when we have time to refactor the client
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
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Position       PositionMessage `json:"position"`
	WallsRemaining int             `json:"wallsRemaining"`
	Ready          bool            `json:"ready"`
	Host           bool            `json:"host"`
}

type PlayerConfigurationMessage struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	PrivatePlayerId string `json:"ppid"`
}

type LobbyMessage struct {
	Players        []PlayerMessage `json:"players"`
	WinnerPlayerId *int            `json:"winnerPlayerId,omitempty"`
}

type LobbyJoin struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type PlayerLeftMessage struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type GameStateStateMessage struct {
	CurrentTurnPlayerId int                  `json:"currentTurnPlayerId"`
	PlayerOne           PlayerMessage        `json:"playerOne"`
	PlayerTwo           PlayerMessage        `json:"playerTwo"`
	Walls               []utils.WallPosition `json:"walls"`
}

func GetPlayerLeftMessage(player player.Player) OMessage {
	return OMessage{
		Type: "playerLeft",
		Payload: PlayerLeftMessage{
			Name: player.Name,
			ID:   int(player.ID),
		},
	}
}

func GetLobbyMessage(players *[]*player.Player) OMessage {
	playersMsg := make([]PlayerMessage, len(*players))
	var winnerPlayerId *int
	for i, p := range *players {
		if p.IsWinner() {
			id := int(p.ID)
			winnerPlayerId = &id
		}
		playersMsg[i] = PlayerMessage{
			ID:    int(p.ID),
			Name:  p.Name,
			Ready: p.Ready,
			Host:  p.Host,
		}
	}
	lobbyMessage := OMessage{
		Type: "lobby",
		Payload: LobbyMessage{
			Players:        playersMsg,
			WinnerPlayerId: winnerPlayerId,
		},
	}

	return lobbyMessage
}

func GetJoinedMessage(player player.Player) OMessage {
	return OMessage{
		Type: "joined",
		Payload: LobbyJoin{
			Name: player.Name,
			ID:   int(player.ID),
		},
	}
}

func GetGameStateMessage(gameState *game.GameState) OMessage {
	var currentTurn int
	var walls []utils.WallPosition

	// Match started
	if gameState.StartTime != nil {
		currentTurn = int(gameState.GetCurrentTurnPlayer().ID)
		walls = gameState.Board.GetWalls()
	}

	p1 := (*gameState.Players)[0]
	p2 := (*gameState.Players)[1]

	gameStateMessage := OMessage{
		Type: "gameState",
		Payload: GameStateStateMessage{
			CurrentTurnPlayerId: currentTurn,
			PlayerOne: PlayerMessage{
				ID:             int(p1.ID),
				Name:           p1.Name,
				Position:       PositionMessage{Row: p1.Position.Row, Col: p1.Position.Column},
				WallsRemaining: p1.WallsRemaining,
				Ready:          p1.Ready,
			},
			PlayerTwo: PlayerMessage{
				ID:             int(p2.ID),
				Name:           p2.Name,
				Position:       PositionMessage{Row: p2.Position.Row, Col: p2.Position.Column},
				WallsRemaining: p2.WallsRemaining,
				Ready:          p2.Ready,
			},
			Walls: walls,
		},
	}

	return gameStateMessage
}

type MatchConfigurationMessage struct {
	Type         string `json:"type"`
	PlayerAmount int    `json:"playerAmount"`
}

func GetMatchConfigurationMessage(playerAmount int) OMessage {
	return OMessage{
		Type: "matchConfiguration",
		Payload: MatchConfigurationMessage{
			Type:         "matchConfiguration",
			PlayerAmount: playerAmount,
		},
	}
}

func GetAlreadyStartedMessage() OMessage {
	return OMessage{
		Type: "alreadyStarted",
	}
}

func GetGameFullMessage() OMessage {
	return OMessage{
		Type: "gameFull",
	}
}

type WillTimeOutMessage struct {
	Span string `json:"span"`
}

func GetWillTimeOutMessage() OMessage {
	return OMessage{
		Type: "willTimeOut",
		Payload: WillTimeOutMessage{
			Span: "10 seconds",
		},
	}
}
