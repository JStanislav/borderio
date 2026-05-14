package gamemanager

import "testing"

func TestAddGame(t *testing.T) {
	games := NewGames()
	gm := NewGameManager(nil)
	games.AddGame("test", gm)

	if games.GetGame("test") == nil {
		t.Error("game manager is nil")
	}
}
