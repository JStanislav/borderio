package messages

import "github.com/JStanislav/quoridor-clone/utils"

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
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	Position       PositionMessage `json:"position"`
	WallsRemaining int             `json:"wallsRemaining"`
}

type GameStateStateMessage struct {
	Type                string               `json:"type"`
	CurrentTurnPlayerId int                  `json:"currentTurnPlayerId"`
	PlayerOne           PlayerMessage        `json:"playerOne"`
	PlayerTwo           PlayerMessage        `json:"playerTwo"`
	Walls               []utils.WallPosition `json:"walls"`
}
