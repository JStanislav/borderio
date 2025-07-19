package player

import "github.com/JStanislav/quoridor-clone/utils"

type PlayerID int

type Player struct {
	ID       PlayerID
	Name     string
	Position utils.GridPosition

	// TODO: Add other player properties
}

func New(name string, position utils.GridPosition) *Player {
	return &Player{
		ID:       1,
		Name:     "Player 1",
		Position: position,
	}
}
