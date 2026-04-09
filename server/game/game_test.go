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
	playerOne := player.New(1, "quoro", utils.GridPosition{Column: 4, Row: 0}, p1StartLine, p1FinishLine)
	playerTwo := player.New(2, "wally", utils.GridPosition{Column: 4, Row: 8}, p2StartLine, p2FinishLine)

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

	gs := NewTwoPlayerMatch([]string{"Player 1", "Player 2"})

	go gs.StartMatch(movesChannel) // use new two player match
	go receiveSelected(movesChannel)

	time.Sleep(1 * time.Second)

	g := gs.Board.(*g.Graph)
	playerOne := gs.Players[0]
	playerTwo := gs.Players[1]

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

func receiveSelected(ch <-chan player.Play) {
	for move := range ch {
		fmt.Println("received", move)
	}
}
