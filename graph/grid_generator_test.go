package graph

import (
	"fmt"
	"testing"
)

func TestGenerateGrid(t *testing.T) {
	graph := New()
	err := graph.GenerateBoard(4, 4, GridPosition{Column: 1, Row: 0}, GridPosition{Column: 2, Row: 3})
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

	graph.PrintGrid(4, 4)

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
	err := graph.GenerateBoard(4, 4, GridPosition{Column: 0, Row: 0}, GridPosition{Column: 3, Row: 3})
	if err != nil {
		t.Error(err)
	}
	graph.PrintGrid(4, 4)
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
	err := graph.GenerateBoard(9, 9, GridPosition{Column: 4, Row: 0}, GridPosition{Column: 4, Row: 8})
	if err != nil {
		panic(err)
	}
	graph.PrintGrid(9, 9)
}

func TestLegalMoves(t *testing.T) {
	graph := generateBasicBoard()
	graph.AddWall(6, 3, 6, 2)
	graph.AddWall(6, 2, 7, 2)
	graph.AddWall(8, 6, 7, 6)
	graph.AddWall(1, 7, 2, 7)
	graph.AddWall(5, 5, 5, 6)
	graph.AddWall(0, 0, 1, 0)

	tests := []struct {
		name                             string
		source, target, opponentPosition GridPosition
		want                             bool
	}{
		{"Illegal", GridPosition{Column: 6, Row: 3}, GridPosition{Column: 6, Row: 2}, graph.PlayerTwoPosition, false},
		{"Illegal", GridPosition{Column: 6, Row: 2}, GridPosition{Column: 7, Row: 2}, graph.PlayerTwoPosition, false},
		{"Illegal", GridPosition{Column: 8, Row: 6}, GridPosition{Column: 7, Row: 6}, graph.PlayerTwoPosition, false},
		{"Illegal", GridPosition{Column: 1, Row: 7}, GridPosition{Column: 2, Row: 7}, graph.PlayerTwoPosition, false},
		{"Illegal", GridPosition{Column: 5, Row: 5}, GridPosition{Column: 5, Row: 6}, graph.PlayerTwoPosition, false},
		{"Illegal", GridPosition{Column: 0, Row: 0}, GridPosition{Column: 1, Row: 0}, graph.PlayerTwoPosition, false},
		{"Legal", GridPosition{Column: 0, Row: 0}, GridPosition{Column: 0, Row: 1}, graph.PlayerTwoPosition, true},
		{"Position occupied", GridPosition{Column: 6, Row: 4}, GridPosition{Column: 5, Row: 4}, GridPosition{Column: 5, Row: 4}, false},
		{"No relation at all", GridPosition{Column: 0, Row: 0}, GridPosition{Column: 5, Row: 4}, GridPosition{Column: 5, Row: 4}, false},
	}

	for _, tt := range tests {

		testname := fmt.Sprintf("%s-[%d,%d]->[%d,%d]", tt.name, tt.source.Column, tt.source.Row, tt.target.Column, tt.target.Row)
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
	err := graph.GenerateBoard(9, 9, GridPosition{Column: 4, Row: 0}, GridPosition{Column: 4, Row: 8})
	if err != nil {
		panic(err)
	}
	return graph
}
