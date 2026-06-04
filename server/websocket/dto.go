package websocket

import (
	"time"

	"github.com/JStanislav/quoridor-clone/gamemanager"
)

type PlayerDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	Host           bool
	WallsRemaining int
}

type GameDTO struct {
	ID          string     `json:"id"`
	PlayerCount int        `json:"playerCount"`
	StartTime   *time.Time `json:"startTime"`
	IsOver      bool
	Players     []PlayerDTO `json:"players"`
	WallLength  int         `json:"wallLength"`

	// Board related
	FinishLineType string `json:"finishLineType"`
	Columns        int    `json:"columns"`
	Rows           int    `json:"rows"`

	CurrentTurn int `json:"currentTurn"`
}

func GetGameDTO(hash string, gm *gamemanager.GameManager) GameDTO {
	players := []PlayerDTO{}
	for _, p := range *gm.Game.Players {
		players = append(players, PlayerDTO{
			ID:             int(p.ID),
			Name:           p.Name,
			Host:           p.Host,
			WallsRemaining: p.WallsRemaining,
		})
	}
	return GameDTO{
		ID:          hash,
		PlayerCount: gm.Game.PlayerCount,
		StartTime:   gm.Game.StartTime,
		IsOver:      gm.IsGameOver(),

		Players:        players,
		WallLength:     gm.Game.WallLength,
		FinishLineType: string(gm.Game.FinishLineType),
		Columns:        gm.Game.Columns,
		Rows:           gm.Game.Rows,
	}
}
