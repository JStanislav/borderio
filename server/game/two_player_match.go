package game

import (
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
)

func NewTwoPlayerMatch(playerNames []string) *GameState {
	columns := 9
	rows := 11

	playerOne := player.New(1, playerNames[0], utils.GridPosition{}, utils.Line{}, utils.Line{})
	playerTwo := player.New(2, playerNames[1], utils.GridPosition{}, utils.Line{}, utils.Line{})

	p1StartPosition := utils.GridPosition{Column: 4, Row: 1}
	p2StartPosition := utils.GridPosition{Column: 4, Row: rows - 2}

	p1StartLine := utils.Line{Type: utils.HorizontalLine, Index: 1}
	p2StartLine := utils.Line{Type: utils.HorizontalLine, Index: rows - 2}
	p1FinishLine := utils.Line{Type: utils.HorizontalLine, Index: rows - 1}
	p2FinishLine := utils.Line{Type: utils.HorizontalLine, Index: 0}

	playerOne.Position = &p1StartPosition
	playerOne.StartLine = p1StartLine
	playerOne.FinishLine = p1FinishLine
	playerTwo.Position = &p2StartPosition
	playerTwo.StartLine = p2StartLine
	playerTwo.FinishLine = p2FinishLine

	g := New(2, []*player.Player{playerOne, playerTwo}, columns, rows, Vertical)
	g.Players = []*player.Player{playerOne, playerTwo}
	g.CurrentTurn = playerOne.ID

	return g
}
