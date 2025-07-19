package game

import (
	"fmt"
	"testing"
	"time"

	"github.com/JStanislav/quoridor-clone/graph"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
)

func TestMain(t *testing.T) {
	graph := graph.New()
	err := graph.GenerateBoard(9, 9, utils.GridPosition{Column: 4, Row: 0}, utils.GridPosition{Column: 4, Row: 8})
	if err != nil {
		panic(err)
	}

	playerOne := player.New("quoro", utils.GridPosition{Column: 4, Row: 0})
	playerTwo := player.New("wally", utils.GridPosition{Column: 4, Row: 8})

	p1Move := utils.GridPosition{Column: 4, Row: 1}
	p2Move := utils.GridPosition{Column: 4, Row: 7}

	if graph.IsLegalMove(playerOne.Position, p1Move, playerTwo.Position) {
		playerOne.Position = p1Move
	}

	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}

	graph.AddWall(4, 1, 4, 2)
	if graph.IsLegalMove(playerOne.Position, utils.GridPosition{Row: 2, Column: 4}, playerTwo.Position) {
		playerOne.Position = p1Move
	}

	p2Move = utils.GridPosition{Row: 6, Column: 4}
	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 5, Column: 4}
	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 4, Column: 4}
	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 3, Column: 4}
	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 2, Column: 4}
	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 2, Column: 3}
	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 1, Column: 3}
	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 1, Column: 5}
	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

	p2Move = utils.GridPosition{Row: 0, Column: 4}
	if graph.IsLegalMove(playerTwo.Position, p2Move, playerOne.Position) {
		playerTwo.Position = p2Move
	}
	graph.PrintGrid(9, 9, playerOne, playerTwo)
	fmt.Println("---------------------------------")

}

func TestMatch(t *testing.T) {
	p1Channel := make(chan utils.GridPosition)
	p2Channel := make(chan utils.GridPosition)
	movesChannel := make(chan PlayerMovement)

	p1StartPosition := utils.GridPosition{Column: 4, Row: 0}
	p2StartPosition := utils.GridPosition{Column: 4, Row: 8}
	playerOne := player.New("quoro", p1StartPosition)
	playerOne.ID = 1
	playerTwo := player.New("wally", p2StartPosition)
	playerTwo.ID = 2

	gs := GameState{}

	go gs.StartMatch(playerOne, playerTwo, p1Channel, p2Channel, movesChannel)
	go receiveSelected(movesChannel)

	time.Sleep(1 * time.Second)

	g := gs.Board.(*graph.Graph)
	g.PrintGrid(9, 9, playerOne, playerTwo)

	p1Channel <- utils.GridPosition{Column: 4, Row: 1}
	p2Channel <- utils.GridPosition{Column: 4, Row: 7}

	time.Sleep(1 * time.Second)
	g.PrintGrid(9, 9, playerOne, playerTwo)
}

func receiveSelected(ch <-chan PlayerMovement) {
	for move := range ch {
		fmt.Println("received", move)
	}
}
