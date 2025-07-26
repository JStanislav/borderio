package graph

import (
	"errors"
	"fmt"

	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
	"github.com/dominikbraun/graph"
)

// TODO: add walls length (default should be 2, in this implementation is 1)
type Board interface {
	GenerateBoard(columns, rows int, playerOneStart, playerTwoStart utils.GridPosition) error
	AddWall(column1, row1, column2, row2 int) error
	IsOccupied(column, row int) (bool, error)
	IsLegalMove(source, target, opponentPosition utils.GridPosition) bool
}

type Graph struct {
	Graph             graph.Graph[string, Cell]
	PlayerOnePosition utils.GridPosition
	PlayerTwoPosition utils.GridPosition
}

func New() *Graph {
	return &Graph{}
}

type Cell struct {
	Id                  string
	IsOccupied          bool
	Column              int
	Row                 int
	IsPlayerOneFinalRow bool
	IsPlayerTwoFinalRow bool
}

func CellHash(c Cell) string {
	return fmt.Sprintf("C%d-R%d", c.Column, c.Row)
}

func (g *Graph) GenerateBoard(columns, rows int, playerOneStart, playerTwoStart utils.GridPosition) error {
	g.Graph = graph.New(CellHash)

	for i := range rows {
		for j := range columns {
			cell := Cell{
				Id:     fmt.Sprintf("R%d-C%d", i, j),
				Column: j,
				Row:    i,
			}
			if i == playerOneStart.Row && j == playerOneStart.Column {
				cell.IsOccupied = true
				g.PlayerOnePosition = playerOneStart
			}
			if i == playerTwoStart.Row && j == playerTwoStart.Column {
				cell.IsOccupied = true
				g.PlayerTwoPosition = playerTwoStart
			}
			g.Graph.AddVertex(cell)

			// Creates the goal line for each player
			if i == 0 {
				cell.IsPlayerTwoFinalRow = true
			}
			if i == rows-1 {
				cell.IsPlayerOneFinalRow = true
			}

		}

	}

	// Creates edges
	for i := range rows - 1 {
		for j := range columns - 1 {
			cell := Cell{
				Column: j,
				Row:    i,
			}
			if err := g.Graph.AddEdge(CellHash(cell), CellHash(Cell{Row: i + 1, Column: j})); err != nil {
				fmt.Println(err)
			}
			if err := g.Graph.AddEdge(CellHash(cell), CellHash(Cell{Row: i, Column: j + 1})); err != nil {
				fmt.Println(err)
			}

			// Literally the border cases
			if i == rows-2 {
				if err := g.Graph.AddEdge(CellHash(Cell{Row: i + 1, Column: j}), CellHash(Cell{Row: i + 1, Column: j + 1})); err != nil {
					fmt.Println(err)
				}
			}
			if j == columns-2 {
				if err := g.Graph.AddEdge(CellHash(Cell{Row: i, Column: j + 1}), CellHash(Cell{Row: i + 1, Column: j + 1})); err != nil {
					fmt.Println(err)
				}
			}

		}
	}

	return nil
}

func (g *Graph) AdjacencyMap() (map[string]map[string]graph.Edge[string], error) {
	return g.Graph.AdjacencyMap()
}

func (g *Graph) AddWall(column1, row1, column2, row2 int) error {
	return g.Graph.RemoveEdge(CellHash(Cell{Column: column1, Row: row1}), CellHash(Cell{Column: column2, Row: row2}))
}

func (g *Graph) IsOccupied(column, row int) (bool, error) {
	cell, err := g.Graph.Vertex(CellHash(Cell{Column: column, Row: row}))
	if err != nil {
		return false, errors.New("vertex not found")
	}
	return cell.IsOccupied, nil
}

func (g *Graph) IsLegalMove(source, target, opponentPosition utils.GridPosition) bool {
	// jugador intenta mover hacia una casilla ocupada
	if target == opponentPosition {
		return false
	}

	// el jugador intenta saltear al otro jugador
	if g.IsAdjacent(source, opponentPosition) && g.IsAdjacent(opponentPosition, target) {
		return true
	}

	return g.IsAdjacent(source, target)
}

func (g *Graph) IsAdjacent(source, target utils.GridPosition) bool {
	_, err := g.Graph.Edge(CellHash(Cell{Column: source.Column, Row: source.Row}), CellHash(Cell{Column: target.Column, Row: target.Row}))
	return !errors.Is(err, graph.ErrEdgeNotFound)
}

func (g *Graph) PrintGrid(columns, rows int, playerOne, playerTwo *player.Player) {
	for i := range rows {
		for j := range columns {
			vertex, err := g.Graph.Vertex(CellHash(Cell{Row: i, Column: j}))

			if err != nil {
				fmt.Printf("Vertex not found: %+v\n", err)
			}
			line := fmt.Sprintf("|%+v|", vertex.Id)
			if playerOne.Position.Row == i && playerOne.Position.Column == j {
				fmt.Printf("\x1b[37;40;%dm%-0s\x1b[37;9;m", 96, line)
			} else if playerTwo.Position.Row == i && playerTwo.Position.Column == j {
				fmt.Printf("\x1b[37;40;%dm%-0s\x1b[37;9;m", 91, line)
			} else {
				fmt.Printf("\x1b[37;40;%dm%-0s\x1b[37;9;m", 97, line)
			}
		}
		fmt.Println()
	}
}
