package game

import (
	"fmt"
	"testing"
	"time"

	"github.com/JStanislav/quoridor-clone/graph"
	g "github.com/JStanislav/quoridor-clone/graph"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
)

func TestMain(t *testing.T) {
	graph := g.New(2, graph.Square)
	err := graph.GenerateBoard(9, 9)
	if err != nil {
		panic(err)
	}

	p1StartLine := utils.Line{Type: utils.HorizontalLine, Index: 0}
	p2StartLine := utils.Line{Type: utils.HorizontalLine, Index: 8}
	p1FinishLine := utils.Line{Type: utils.HorizontalLine, Index: 8}
	p2FinishLine := utils.Line{Type: utils.HorizontalLine, Index: 0}
	playerOne := player.New(1, "ppid1", "quoro", utils.GridPosition{Column: 4, Row: 0}, 9, p1StartLine, p1FinishLine)
	playerTwo := player.New(2, "ppid2", "wally", utils.GridPosition{Column: 4, Row: 8}, 9, p2StartLine, p2FinishLine)

	p1Move := utils.GridPosition{Column: 4, Row: 1}
	p2Move := utils.GridPosition{Column: 4, Row: 7}

	if graph.IsLegalMove(*playerOne.Position, p1Move, []*utils.GridPosition{playerTwo.Position}) {
		playerOne.Position = &p1Move
	}

	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}

	graph.AddWall(g.Horizontal, utils.WallPosition{CellA: utils.GridPosition{Column: 4, Row: 1}, CellB: utils.GridPosition{Column: 4, Row: 2}}) //4, 1, 4, 2
	if graph.IsLegalMove(*playerOne.Position, utils.GridPosition{Row: 2, Column: 4}, []*utils.GridPosition{playerTwo.Position}) {
		playerOne.Position = &p1Move
	}

	p2Move = utils.GridPosition{Row: 6, Column: 4}
	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 5, Column: 4}
	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 4, Column: 4}
	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 3, Column: 4}
	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 2, Column: 4}
	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 2, Column: 3}
	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 1, Column: 3}
	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 1, Column: 5}
	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 0, Column: 4}
	if graph.IsLegalMove(*playerTwo.Position, p2Move, []*utils.GridPosition{playerOne.Position}) {
		playerTwo.Position = &p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

}

func TestMatch(t *testing.T) {
	movesChannel := make(chan player.Play)

	gs := NewTwoPlayerMatch()
	p1 := player.New(1, "ppid1", "quoro", utils.GridPosition{Column: 4, Row: 0}, 9, utils.Line{Type: utils.HorizontalLine, Index: 0}, utils.Line{Type: utils.HorizontalLine, Index: 8})
	p2 := player.New(2, "ppid2", "wally", utils.GridPosition{Column: 4, Row: 8}, 9, utils.Line{Type: utils.HorizontalLine, Index: 8}, utils.Line{Type: utils.HorizontalLine, Index: 0})

	gs.AddPlayer(p1)
	gs.AddPlayer(p2)

	gs.StartMatch(movesChannel) // use new two player match
	go receiveSelected(movesChannel)

	time.Sleep(1 * time.Second)

	g := gs.Board.(*g.Graph)
	players := *gs.Players
	playerOne := players[0]
	playerTwo := players[1]

	g.PrintGrid(9, 9, playerOne, playerTwo)

	plays := []struct {
		player player.Player
		play   player.Play
	}{
		{*playerOne, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 1}}},
		{*playerTwo, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 7}}},
		{*playerTwo, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 6}}},
		{*playerTwo, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 6}}},
		{*playerTwo, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 6}}},
		{*playerTwo, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 6}}},
		{*playerOne, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 2}}},
		{*playerTwo, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 6}}},
		{*playerOne, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 3}}},
		{*playerTwo, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 5}}},
		{*playerOne, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 4}}},
		{*playerTwo, player.Play{PlayType: player.PlayerMove, Position: &utils.GridPosition{Column: 4, Row: 3}}},
	}

	for _, play := range plays {
		if err := play.player.OnPlayerPlay(play.player.ID, play.play); err != nil {
			fmt.Println("err", err)
		}
		time.Sleep(1500 * time.Millisecond)
		g.PrintGrid(9, 9, playerOne, playerTwo)
	}
}

func TestGameHosts(t *testing.T) {
	gs := NewTwoPlayerMatch()
	p1 := player.New(1, "ppid1", "quoro", utils.GridPosition{Column: 4, Row: 0}, 9, utils.Line{Type: utils.HorizontalLine, Index: 0}, utils.Line{Type: utils.HorizontalLine, Index: 8})
	p2 := player.New(2, "ppid2", "wally", utils.GridPosition{Column: 4, Row: 8}, 9, utils.Line{Type: utils.HorizontalLine, Index: 8}, utils.Line{Type: utils.HorizontalLine, Index: 0})

	gs.AddPlayer(p1)
	gs.AddPlayer(p2)

	if p1.Host != true {
		t.Errorf("expected player 1 to be host")
	}

	gs.RemovePlayer(1)

	if p2.Host != true {
		t.Errorf("expected player 2 to be host after player 1 leaves")
	}

}

func receiveSelected(ch <-chan player.Play) {
	for move := range ch {
		fmt.Println("received", move)
	}
}
