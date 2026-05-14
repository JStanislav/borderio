package messages

import "github.com/JStanislav/quoridor-clone/utils"

type Message[T any] struct {
	Type       string            `json:"type"`
	PlayerId   int               `json:"playerId"`
	Target     PositionMessage   `json:"target"`     // remove this and put it in Payload when we have time to refactor the client
	WallTarget WallTargetMessage `json:"wallTarget"` // remove this and put it in Payload when we have time to refactor the client
	Payload    T                 `json:"payload"`

	PrivatePlayerId string `json:"ppid"`
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
}

type PlayerConfigurationMessage struct {
	Type            string `json:"type"`
	ID              int    `json:"id"`
	Name            string `json:"name"`
	PrivatePlayerId string `json:"ppid"`
}

type LobbyMessage struct {
	Type    string          `json:"type"`
	Players []PlayerMessage `json:"players"`
}

type GameStateStateMessage struct {
	Type                string               `json:"type"`
	CurrentTurnPlayerId int                  `json:"currentTurnPlayerId"`
	PlayerOne           PlayerMessage        `json:"playerOne"`
	PlayerTwo           PlayerMessage        `json:"playerTwo"`
	Walls               []utils.WallPosition `json:"walls"`
}
