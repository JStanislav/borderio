package game

import (
	"time"

	"github.com/JStanislav/quoridor-clone/player"
)

type GameStats struct {
	PlayerWinnerId      int
	TotalMoves          int
	TotalWallPlacements int
	TotalWallsRemaining int
	StartTime           *time.Time
	EndTime             *time.Time
	Players             *[]*player.Player

	Points int

	Steps []player.Play
}
