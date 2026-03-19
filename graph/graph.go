package graph

import (
	"errors"
	"fmt"
	"slices"

	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
	"github.com/dominikbraun/graph"
)

type WallType string

const (
	Horizontal WallType = "horizontal"
	Vertical   WallType = "vertical"
	Undefined  WallType = "undefined"
)

type Board interface {
	GenerateBoard(columns, rows int, playerOneStart, playerTwoStart utils.GridPosition) error
	AddWall(wallType WallType, start utils.WallPosition) error
	IsOccupied(column, row int) (bool, error)
	IsLegalMove(source, target, opponentPosition utils.GridPosition) bool
	IsLegalMove2(source, target utils.GridPosition) bool
	GetWalls() []utils.WallPosition
}

type Graph struct {
	Graph             graph.Graph[string, Cell]
	PlayerOnePosition utils.GridPosition
	PlayerTwoPosition utils.GridPosition
	Walls             []utils.WallPosition
	wallLength        int
}

func New(wallLength int) *Graph {
	return &Graph{wallLength: wallLength}
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

func (g *Graph) IsWallOccupied(position utils.WallPosition) bool {
	hashA := CellHash(Cell{Column: position.CellA.Column, Row: position.CellA.Row})
	hashB := CellHash(Cell{Column: position.CellB.Column, Row: position.CellB.Row})
	_, err := g.Graph.Edge(hashA, hashB)
	return errors.Is(err, graph.ErrEdgeNotFound)
}

func (g *Graph) AddWall(wallType WallType, start utils.WallPosition) error {
	_g, err := g.Graph.Clone()
	if err != nil {
		return err
	}
	var errCreatingWall error

	if wallType == Undefined {
		// TODO: remove this crap
		if start.CellA.Row == start.CellB.Row || start.CellA.Column != start.CellB.Column {
			wallType = Vertical
		}
		if start.CellA.Column == start.CellB.Column || start.CellA.Row != start.CellB.Row {
			wallType = Horizontal
		}
	}

	if wallType == Horizontal {
		if start.CellA.Row == start.CellB.Row || start.CellA.Column != start.CellB.Column {
			return errors.New("wall is not horizontal")
		}
		end := start.CellA.Column + g.wallLength
		for j := start.CellA.Column; j < end; j++ {

			if j != end-1 {
				// check if a wall is cut through

				if slices.Contains(g.GetWalls(), utils.WallPosition{CellA: utils.GridPosition{Column: j, Row: start.CellA.Row}, CellB: utils.GridPosition{Column: j + 1, Row: start.CellA.Row}}) {
					errCreatingWall = errors.New("wall is cut through")
					str := fmt.Sprintf("j: %d, wallOccupied top: %t", j, g.IsWallOccupied(utils.WallPosition{CellA: utils.GridPosition{Column: j, Row: start.CellA.Row}, CellB: utils.GridPosition{Column: j + 1, Row: start.CellA.Row}}))
					str2 := fmt.Sprintf("j: %d, wallOccupied bottom: %t", j, g.IsWallOccupied(utils.WallPosition{CellA: utils.GridPosition{Column: j, Row: start.CellB.Row}, CellB: utils.GridPosition{Column: j + 1, Row: start.CellB.Row}}))
					fmt.Println(str)
					fmt.Println(str2)

				}
			}

			err := _g.RemoveEdge(CellHash(Cell{Column: j, Row: start.CellA.Row}), CellHash(Cell{Column: j, Row: start.CellB.Row}))
			if err != nil {
				errCreatingWall = err
			}
		}

		if errCreatingWall != nil {
			return fmt.Errorf("error creating wall: %s", errCreatingWall)
		}

	}

	if wallType == Vertical {
		if start.CellA.Column == start.CellB.Column || start.CellA.Row != start.CellB.Row {
			return errors.New("wall is not vertical")
		}
		end := start.CellA.Row + g.wallLength
		for i := start.CellA.Row; i < end; i++ {

			if i != end-1 {
				// check if a wall is cut through
				if slices.Contains(g.GetWalls(), utils.WallPosition{CellA: utils.GridPosition{Column: start.CellA.Column, Row: i}, CellB: utils.GridPosition{Column: start.CellA.Column, Row: i + 1}}) {
					errCreatingWall = errors.New("wall is cut through")
					str := fmt.Sprintf("i: %d, wallOccupied top: %t", i, g.IsWallOccupied(utils.WallPosition{CellA: utils.GridPosition{Row: i, Column: start.CellA.Column}, CellB: utils.GridPosition{Row: i + 1, Column: start.CellA.Column}}))
					str2 := fmt.Sprintf("i: %d, wallOccupied bottom: %t", i, g.IsWallOccupied(utils.WallPosition{CellA: utils.GridPosition{Row: i, Column: start.CellB.Column}, CellB: utils.GridPosition{Row: i + 1, Column: start.CellB.Column}}))
					fmt.Println(str)
					fmt.Println(str2)

				}
			}

			err := _g.RemoveEdge(CellHash(Cell{Row: i, Column: start.CellA.Column}), CellHash(Cell{Row: i, Column: start.CellB.Column}))
			if err != nil {
				errCreatingWall = err
			}
		}

		if errCreatingWall != nil {
			return fmt.Errorf("error creating wall: %s", errCreatingWall)
		}

	}

	g.Walls = append(g.Walls, utils.WallPosition{CellA: start.CellA, CellB: start.CellB})
	g.Graph = _g
	return nil
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

func (g *Graph) IsLegalMove2(source, target utils.GridPosition) bool {
	occupied, err := g.IsOccupied(target.Column, target.Row)
	if err != nil {
		return false
	}

	if !occupied {
		return g.IsAdjacent(source, target)
	}

	return false
}

func (g *Graph) GetWalls() []utils.WallPosition {
	// Implementation for getting walls
	return g.Walls
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
