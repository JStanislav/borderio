package game

import (
	"errors"

	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
)

type TwoPlayerMatch struct {
	GameState
}

func NewTwoPlayerMatch() *TwoPlayerMatch {
	columns := 9
	rows := 11

	g := New(2, 2, columns, rows, Vertical)

	return &TwoPlayerMatch{GameState: *g}
}

func (m *TwoPlayerMatch) AddPlayer(p *player.Player) error {
	id := m.GameState.GetUnusedPlayerID()
	if id > 2 {
		return errors.New("game is full")
	}
	p.ID = player.PlayerID(id)

	if p.ID == 1 {
		p1StartPosition := utils.GridPosition{Column: 4, Row: 1}
		p1StartLine := utils.Line{Type: utils.HorizontalLine, Index: 1}
		p1FinishLine := utils.Line{Type: utils.HorizontalLine, Index: m.Rows - 1}

		p.Position = &p1StartPosition
		p.StartLine = p1StartLine
		p.FinishLine = p1FinishLine
	}

	if p.ID == 2 {
		p2StartPosition := utils.GridPosition{Column: 4, Row: m.Rows - 2}
		p2StartLine := utils.Line{Type: utils.HorizontalLine, Index: m.Rows - 2}
		p2FinishLine := utils.Line{Type: utils.HorizontalLine, Index: 0}

		p.Position = &p2StartPosition
		p.StartLine = p2StartLine
		p.FinishLine = p2FinishLine

	}

	return m.GameState.AddPlayer(p)
}
