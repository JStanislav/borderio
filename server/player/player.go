package player

import (
	"github.com/JStanislav/quoridor-clone/utils"
)

type PlayerID int
type PlayType int

const (
	InvalidPlayID PlayType = -1
	PlayerMove             = iota + 1
	WallPlacement
)

type Play struct {
	PlayType   PlayType
	Position   *utils.GridPosition
	WallPlaced *utils.WallPosition
}

type OnPlayerPlay func(playerID PlayerID, play Play) error

type Player struct {
	ID       PlayerID
	Name     string
	Position *utils.GridPosition

	OnPlayerPlay OnPlayerPlay

	StartLine  utils.Line
	FinishLine utils.Line
}

func New(id PlayerID, name string, position utils.GridPosition, startLine utils.Line, finishLine utils.Line) *Player {
	return &Player{
		ID:       id,
		Name:     name,
		Position: &position,

		StartLine:  startLine,
		FinishLine: finishLine,
	}
}

func GetPlayersPositions(players []*Player) []*utils.GridPosition {
	positions := make([]*utils.GridPosition, len(players))
	for i, p := range players {
		positions[i] = &utils.GridPosition{
			Row:    p.Position.Row,
			Column: p.Position.Column,
		}
	}
	return positions
}

func GetPlayersFinishLines(players []*Player) []utils.Line {
	finishLines := make([]utils.Line, len(players))
	for i, p := range players {
		finishLines[i] = p.FinishLine
	}
	return finishLines
}

func (p *Player) IsWinner() bool {
	if p.Position == nil {
		return false
	}
	if p.FinishLine.Type == utils.HorizontalLine {
		return p.Position.Row == p.FinishLine.Index
	} else {
		return p.Position.Column == p.FinishLine.Index
	}
}

func (p *Player) IsFinishLine(play Play) bool {
	if p.Position == nil {
		return false
	}
	if p.FinishLine.Type == utils.HorizontalLine {
		return play.Position.Row == p.FinishLine.Index
	} else {
		return play.Position.Column == p.FinishLine.Index
	}
}
