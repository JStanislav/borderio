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

	g := New(2, 2, columns, rows, Horizontal)

	return &TwoPlayerMatch{GameState: *g}
}

func (m *TwoPlayerMatch) AddPlayer(name string, ppid string, walls int) error {
	id := m.GameState.GetUnusedPlayerID()
	if id > 2 {
		return errors.New("game is full")
	}

	var startPosition utils.GridPosition
	var startLine utils.Line
	var finishLine utils.Line
	if id == 1 {
		startPosition = utils.GridPosition{Column: 4, Row: 1}
		startLine = utils.Line{Type: utils.HorizontalLine, Index: 1}
		finishLine = utils.Line{Type: utils.HorizontalLine, Index: m.Rows - 1}
	}

	if id == 2 {
		startPosition = utils.GridPosition{Column: 4, Row: m.Rows - 2}
		startLine = utils.Line{Type: utils.HorizontalLine, Index: m.Rows - 2}
		finishLine = utils.Line{Type: utils.HorizontalLine, Index: 0}
	}

	p := player.New(ppid, name, startPosition, walls, startLine, finishLine)
	p.ID = player.PlayerID(id)

	return m.GameState.AddPlayer(p)
}
