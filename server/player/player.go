package player

import "github.com/JStanislav/quoridor-clone/utils"

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
	// TODO: Add other player properties
}

func New(name string, position utils.GridPosition) *Player {
	return &Player{
		ID:       1,
		Name:     "Player 1",
		Position: &position,
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
