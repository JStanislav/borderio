package gamemanager

import (
	"context"
	"testing"
	"time"

	"github.com/JStanislav/quoridor-clone/game"
)

func TestAddGame(t *testing.T) {
	games := NewGames()

	gs := game.New(2, 2, 8, 8, game.Horizontal)

	gm := NewGameManager(context.Background(), "id-test", gs, nil, 30*time.Second)
	games.AddGame("test", gm)

	if games.GetGame("test") == nil {
		t.Error("game manager is nil")
	}
}

func newStartedTestGame(t *testing.T) (*GameManager, *IO, *IO) {
	t.Helper()

	match := game.NewTwoPlayerMatch()
	gm := NewGameManager(context.Background(), "id-test", &match.GameState, func(game.GameStats) error { return nil }, 30*time.Second)
	gm.ReconnectionGracePeriod = 50 * time.Millisecond

	for _, ppid := range []string{"ppid1", "ppid2"} {
		if err := match.AddPlayer(ppid, ppid, 8); err != nil {
			t.Fatal(err)
		}
	}

	io1 := NewIO("ppid1", "ppid1", nil)
	io2 := NewIO("ppid2", "ppid2", nil)
	gm.handleJoin(io1)
	gm.handleJoin(io2)

	gm.Game.StartMatchWithMovementsChannel()
	gm.Started = true

	return gm, io1, io2
}

func expireGracePeriod(t *testing.T, gm *GameManager) {
	t.Helper()

	select {
	case ppid := <-gm.graceExpired:
		gm.handleGraceExpired(ppid)
	case <-time.After(time.Second):
		t.Fatal("grace period did not expire")
	}
}

func TestOpponentWinsWhenGracePeriodExpires(t *testing.T) {
	gm, io1, _ := newStartedTestGame(t)

	gm.handleLeave(io1)

	if len(*gm.Game.Players) != 2 || gm.Lobbyists["ppid1"] == nil {
		t.Fatalf("player should be kept in the game during the grace period")
	}
	if gm.Lobbyists["ppid1"].Connected {
		t.Fatalf("player should be marked as disconnected")
	}

	expireGracePeriod(t, gm)

	if !gm.IsGameOver() {
		t.Fatal("game should be over")
	}
	if winner := gm.Game.GetWinnerPlayerId(); winner == nil || *winner != 2 {
		t.Fatalf("expected player 2 to win, got %v", winner)
	}
}

func TestReconnectionWithinGracePeriod(t *testing.T) {
	gm, io1, _ := newStartedTestGame(t)

	gm.handleLeave(io1)
	gm.handleJoin(NewIO("ppid1", "ppid1", nil))

	if len(gm.disconnected) != 0 || !gm.Lobbyists["ppid1"].Connected {
		t.Fatal("player should be connected again")
	}
	if gm.Lobbyists["ppid1"].Spectator {
		t.Fatal("player should not come back as a spectator")
	}

	select {
	case <-gm.graceExpired:
		t.Fatal("grace period timer should have been stopped")
	case <-time.After(100 * time.Millisecond):
	}

	if gm.IsGameOver() || gm.Game.GetWinnerPlayerId() != nil {
		t.Fatal("game should still be running")
	}
}

func TestDrawWhenNoPlayerComesBack(t *testing.T) {
	gm, io1, io2 := newStartedTestGame(t)

	// a connected spectator doesn't prevent the draw
	gm.handleJoin(NewIO("spectator", "spectator", nil))

	gm.handleLeave(io1)
	gm.handleLeave(io2)

	expireGracePeriod(t, gm)

	if !gm.IsGameOver() {
		t.Fatal("game should be over")
	}
	if !gm.Game.Draw || gm.Game.GetWinnerPlayerId() != nil {
		t.Fatalf("expected a draw, got winner %v", gm.Game.GetWinnerPlayerId())
	}
	if gm.IsAwaitingReconnection() {
		t.Fatal("no reconnection should be pending after the game ended")
	}
}

func TestSpectatorLeavingDoesNotStartGracePeriod(t *testing.T) {
	gm, _, _ := newStartedTestGame(t)

	spectator := NewIO("spectator", "spectator", nil)
	gm.handleJoin(spectator)
	gm.handleLeave(spectator)

	if len(gm.disconnected) != 0 {
		t.Fatal("spectators should not get a grace period")
	}
	if _, ok := gm.Lobbyists["spectator"]; ok {
		t.Fatal("spectator should be removed from the lobby")
	}
}

func TestGCKeepsGamesAwaitingReconnection(t *testing.T) {
	gm, io1, io2 := newStartedTestGame(t)
	games := NewGames()
	games.AddGame("test", gm)

	gm.handleLeave(io1)
	gm.handleLeave(io2)

	games.DeleteOldGames()
	if games.GetGame("test") == nil {
		t.Fatal("game awaiting reconnection should not be deleted")
	}

	gm.handleJoin(NewIO("ppid1", "ppid1", nil))
	expireGracePeriod(t, gm)
	if !gm.IsGameOver() || gm.Game.Draw {
		t.Fatal("the connected player should win")
	}
}
