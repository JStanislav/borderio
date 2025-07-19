package graph

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
)

func TestGenerateGrid(t *testing.T) {
	graph := New()
	err := graph.GenerateBoard(4, 4, utils.GridPosition{Column: 1, Row: 0}, utils.GridPosition{Column: 2, Row: 3})
	if err != nil {
		t.Error(err)
	}
	adjacencyMap, err := graph.AdjacencyMap()
	if err != nil {
		t.Error(err)
	}

	var tests = []struct {
		col, row int
		want     int
	}{
		{0, 0, 2},
		{0, 3, 2},
		{2, 2, 4},
		{0, 1, 3},
		{3, 3, 2},
	}

	p1, p2 := player.Player{}, player.Player{}
	graph.PrintGrid(4, 4, &p1, &p2)

	for _, tt := range tests {

		testname := fmt.Sprintf("%d,%d", tt.col, tt.row)
		t.Run(testname, func(t *testing.T) {
			node := adjacencyMap[CellHash(Cell{Column: tt.col, Row: tt.row})]
			ans := len(node)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestAddWall(t *testing.T) {
	graph := New()
	err := graph.GenerateBoard(4, 4, utils.GridPosition{Column: 0, Row: 0}, utils.GridPosition{Column: 3, Row: 3})
	if err != nil {
		t.Error(err)
	}
	p1 := player.Player{}
	p2 := player.Player{}
	graph.PrintGrid(4, 4, &p1, &p2)
	err = graph.AddWall(0, 0, 0, 1)
	if err != nil {
		t.Error(err)
	}
	err = graph.AddWall(0, 0, 1, 0)
	if err != nil {
		t.Error(err)
	}
	adjacencyMap, _ := graph.AdjacencyMap()

	fmt.Println("Added wall", adjacencyMap[CellHash(Cell{Column: 0, Row: 0})])
}

func TestPrintGrid(t *testing.T) {
	graph := New()
	err := graph.GenerateBoard(9, 9, utils.GridPosition{Column: 4, Row: 0}, utils.GridPosition{Column: 4, Row: 8})
	if err != nil {
		panic(err)
	}
	p1 := player.Player{}
	p2 := player.Player{}
	graph.PrintGrid(9, 9, &p1, &p2)
}

func TestLegalMoves(t *testing.T) {
	graph := generateBasicBoard()
	graph.AddWall(6, 3, 6, 2)
	graph.AddWall(6, 2, 7, 2)
	graph.AddWall(8, 6, 7, 6)
	graph.AddWall(1, 7, 2, 7)
	graph.AddWall(5, 5, 5, 6)
	graph.AddWall(0, 0, 1, 0)
	/*
		|R0-C0░R0-C1░R0-C2░|R0-C3░R0-C4░R0-C5░R0-C6░R0-C7░R0-C8|
		░░░░░░█░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
		|R1-C0░R1-C1░R1-C2░|R1-C3░R1-C4░R1-C5░R1-C6░R1-C7░R1-C8|
		░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
		|R2-C0░R2-C1░R2-C2░|R2-C3░R2-C4░R2-C5░R2-C6░R2-C7░R2-C8|
		░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░█████░█████░░░░░░░
		|R3-C0░R3-C1░R3-C2░|R3-C3░R3-C4░R3-C5░R3-C6░R3-C7░R3-C8|
		░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
		|R4-C0░R4-C1░R4-C2░|R4-C3░R4-C4░R4-C5░R4-C6░R4-C7░R4-C8|
		░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
		|R5-C0░R5-C1░R5-C2░|R5-C3░R5-C4░R5-C5░R5-C6░R5-C7░R5-C8|
		░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░█████░░░░░░░░░░░░░░░░░░░
		|R6-C0░R6-C1░R6-C2░|R6-C3░R6-C4░R6-C5░R6-C6░R6-C7█R6-C8|
		░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
		|R7-C0░R7-C1█R7-C2░|R7-C3░R7-C4░R7-C5░R7-C6░R7-C7░R7-C8|
		░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
		|R8-C0░R8-C1░R8-C2░|R8-C3░R8-C4░R8-C5░R8-C6░R8-C7░R8-C8|
		█
	*/

	// a little helper so its more readable the test cases
	// converts hashes to grid positions
	// i.e. "R3-C6" -> utils.GridPosition{Column: 6, Row: 3}
	p := func(hash string) utils.GridPosition {
		strs := strings.Split(hash, "-")
		strs[0], _ = strings.CutPrefix(strs[0], "R")
		strs[1], _ = strings.CutPrefix(strs[1], "C")
		c, _ := strconv.Atoi(strs[1])
		r, _ := strconv.Atoi(strs[0])
		return utils.GridPosition{Column: c, Row: r}
	}

	tests := []struct {
		name                             string
		source, target, opponentPosition utils.GridPosition
		want                             bool
	}{
		{"Illegal", p("R3-C6"), p("R2-C6"), graph.PlayerTwoPosition, false},
		{"Illegal", p("R2-C6"), p("R2-C7"), graph.PlayerTwoPosition, false},
		{"Illegal", p("R6-C8"), p("R6-C7"), graph.PlayerTwoPosition, false},
		{"Illegal", p("R7-C1"), p("R7-C2"), graph.PlayerTwoPosition, false},
		{"Illegal", p("R5-C5"), p("R6-C5"), graph.PlayerTwoPosition, false},
		{"Illegal", p("R0-C0"), p("R0-C1"), graph.PlayerTwoPosition, false},
		{"Legal", p("R0-C0"), p("R1-C0"), graph.PlayerTwoPosition, true},
		{"Position occupied", p("R4-C6"), p("R4-C5"), p("R4-C5"), false},
		{"No relation at all", p("R0-C0"), p("R4-C5"), p("R4-C5"), false},
		{"Simple skip", p("R2-C2"), p("R2-C4"), p("R2-C3"), true},
		{"Skip with wall", p("R3-C5"), p("R4-C4"), p("R4-C5"), true},
		{"Illegal skip through wall", p("R5-C7"), p("R5-C5"), p("R6-C5"), false},
	}

	for _, tt := range tests {

		testname := fmt.Sprintf("%s-[R%d,C%d]->[R%d,C%d]", tt.name, tt.source.Row, tt.source.Column, tt.target.Row, tt.target.Column)
		t.Run(testname, func(t *testing.T) {
			isValid := graph.IsLegalMove(tt.source, tt.target, tt.opponentPosition)
			if isValid != tt.want {
				t.Errorf("got %t, want %t", isValid, tt.want)
			}
		})
	}
}

func generateBasicBoard() *Graph {
	graph := New()
	err := graph.GenerateBoard(9, 9, utils.GridPosition{Column: 4, Row: 0}, utils.GridPosition{Column: 4, Row: 8})
	if err != nil {
		panic(err)
	}
	return graph
}
