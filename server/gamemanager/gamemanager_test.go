package gamemanager

import (
	"testing"
	"time"

	"github.com/JStanislav/quoridor-clone/game"
)

func TestAddGame(t *testing.T) {
	games := NewGames()

	gs := game.New(2, 2, 8, 8, game.Horizontal)

	gm := NewGameManager(gs, nil, 30*time.Second)
	games.AddGame("test", gm)

	if games.GetGame("test") == nil {
		t.Error("game manager is nil")
	}
}
