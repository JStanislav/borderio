package game

import (
	"time"

	"github.com/JStanislav/quoridor-clone/graph"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
)

type GameState struct {
	Board       graph.Board
	StartTime   *time.Time
	CurrentTurn player.PlayerID
}

func StartMatch() {
	board := graph.New()
	p1StartPosition := utils.GridPosition{Column: 4, Row: 0}
	p2StartPosition := utils.GridPosition{Column: 4, Row: 8}

	board.GenerateBoard(9, 9, p1StartPosition, p2StartPosition)

	playerOne := player.New("quoro", p1StartPosition)
	playerTwo := player.New("wally", p2StartPosition)

	p1Move := utils.GridPosition{Column: 4, Row: 1}
	p2Move := utils.GridPosition{Column: 4, Row: 7}

	if board.IsLegalMove(playerOne.Position, p1Move, playerTwo.Position) {
		playerOne.Position = p1Move
	}

	if board.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p1Move
	}
}
