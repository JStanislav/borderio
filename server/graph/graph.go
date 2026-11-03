package graph

import (
	"errors"
	"fmt"
	"slices"

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
	GenerateBoard(columns, rows int) error
	AddWall(wallType WallType, start utils.WallPosition) error
	IsLegalMove(source, target utils.GridPosition, opponentPosition []*utils.GridPosition) bool
	RemoveWall(wallType WallType, start utils.WallPosition) error
	ExistsPath(source, target utils.GridPosition) bool
	GetWalls() []utils.WallPosition
}

type Graph struct {
	Graph      graph.Graph[string, Cell]
	Walls      []utils.WallPosition
	wallLength int
	boardType  BoardType
}

func New(wallLength int, boardType BoardType) *Graph {
	return &Graph{
		wallLength: wallLength,
		boardType:  boardType,
		Graph:      nil,
		Walls:      *new([]utils.WallPosition),
	}
}

type BoardType string

const (
	ExtraRows    BoardType = "extra-rows"
	ExtraColumns BoardType = "extra-columns"
	Square       BoardType = "square"
)

type Cell struct {
	Id     string
	Column int
	Row    int
}

func CellHash(c Cell) string {
	return fmt.Sprintf("C%d-R%d", c.Column, c.Row)
}

func (g *Graph) GenerateBoard(columns, rows int) error {
	g.Graph = graph.New(CellHash)

	for i := range rows {
		for j := range columns {
			cell := Cell{
				Id:     fmt.Sprintf("R%d-C%d", i, j),
				Column: j,
				Row:    i,
			}
			g.Graph.AddVertex(cell)
		}

	}

	// Creates edges
	for i := range rows - 1 {
		for j := range columns - 1 {
			cell := Cell{
				Column: j,
				Row:    i,
			}
			g.Graph.AddEdge(CellHash(cell), CellHash(Cell{Row: i + 1, Column: j}))
			g.Graph.AddEdge(CellHash(cell), CellHash(Cell{Row: i, Column: j + 1}))

			// Literally the border cases
			if i == rows-2 {
				g.Graph.AddEdge(CellHash(Cell{Row: i + 1, Column: j}), CellHash(Cell{Row: i + 1, Column: j + 1}))
			}
			if j == columns-2 {
				g.Graph.AddEdge(CellHash(Cell{Row: i, Column: j + 1}), CellHash(Cell{Row: i + 1, Column: j + 1}))
			}

		}
	}

	removeEdgesFromFinishLines(g.Graph, columns, rows, g.boardType)

	return nil
}

func removeEdgesFromFinishLines(g graph.Graph[string, Cell], columns, rows int, boardType BoardType) {
	switch boardType {
	case ExtraRows:
		for i := 0; i < columns-1; i++ {
			g.RemoveEdge(CellHash(Cell{Row: 0, Column: i}), CellHash(Cell{Row: 0, Column: i + 1}))
			g.RemoveEdge(CellHash(Cell{Row: rows - 1, Column: i}), CellHash(Cell{Row: rows - 1, Column: i + 1}))
		}
	case ExtraColumns:
		for i := 0; i < rows-1; i++ {
			g.RemoveEdge(CellHash(Cell{Row: i, Column: 0}), CellHash(Cell{Row: i + 1, Column: 0}))
			g.RemoveEdge(CellHash(Cell{Row: i, Column: columns - 1}), CellHash(Cell{Row: i + 1, Column: columns - 1}))
		}
	}
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

	if wallType == Undefined {
		wallType = infereWallType(start)
	}

	horizontal := wallType == Horizontal

	if horizontal {
		if start.CellA.Row > start.CellB.Row {
			start = utils.WallPosition{CellA: start.CellB, CellB: start.CellA}
		}

		for i := 0; i < g.wallLength-1; i++ {
			for j := 0; j < g.wallLength-1; j++ {
				cellA := utils.GridPosition{Column: start.CellA.Column + j, Row: start.CellA.Row - i}
				cellB := utils.GridPosition{Column: start.CellA.Column + j + 1, Row: start.CellA.Row - i}
				completeWallExists := slices.Contains(g.Walls, utils.WallPosition{CellA: cellA, CellB: cellB}) || slices.Contains(g.Walls, utils.WallPosition{CellA: cellB, CellB: cellA})
				if completeWallExists {
					return errors.New("wall is cut through another wall")
				}
			}
		}
		for i := 0; i < g.wallLength; i++ {
			if err := _g.RemoveEdge(CellHash(Cell{Column: start.CellA.Column + i, Row: start.CellA.Row}), CellHash(Cell{Column: start.CellB.Column + i, Row: start.CellB.Row})); err != nil {
				err = fmt.Errorf("Error removing edge between %+v and %+v: %s\n", start.CellA, start.CellB, err)
				return err
			}
		}
	} else {
		if start.CellA.Column > start.CellB.Column {
			start = utils.WallPosition{CellA: start.CellB, CellB: start.CellA}
		}

		for i := 0; i < g.wallLength-1; i++ {
			for j := 0; j < g.wallLength-1; j++ {
				cellA := utils.GridPosition{Column: start.CellA.Column - i, Row: start.CellA.Row + j}
				cellB := utils.GridPosition{Column: start.CellA.Column - i, Row: start.CellA.Row + j + 1}
				completeWallExists := slices.Contains(g.Walls, utils.WallPosition{CellA: cellA, CellB: cellB}) || slices.Contains(g.Walls, utils.WallPosition{CellA: cellB, CellB: cellA})
				if completeWallExists {
					return errors.New("wall is cut through another wall")
				}
			}
		}
		for i := 0; i < g.wallLength; i++ {
			cellHashA := CellHash(Cell{Column: start.CellA.Column, Row: start.CellA.Row + i})
			cellHashB := CellHash(Cell{Column: start.CellB.Column, Row: start.CellB.Row + i})

			if err := _g.RemoveEdge(cellHashA, cellHashB); err != nil {
				err = fmt.Errorf("Error removing edge between %+v and %+v: %s\n", start.CellA, start.CellB, err)
				return err
			}
		}
	}

	g.Walls = append(g.Walls, utils.WallPosition{CellA: start.CellA, CellB: start.CellB})
	g.Graph = _g
	return nil
}

func (g *Graph) RemoveWall(wallType WallType, start utils.WallPosition) error {
	s := slices.DeleteFunc(g.Walls, func(wall utils.WallPosition) bool {
		return (wall.CellA == start.CellA && wall.CellB == start.CellB) || (wall.CellA == start.CellB && wall.CellB == start.CellA)
	})

	if len(s) < len(g.Walls) {
		g.Walls = s
	}

	_g, err := g.Graph.Clone()
	if err != nil {
		return err
	}

	if wallType == Undefined {
		wallType = infereWallType(start)
	}

	if wallType == Horizontal {
		if start.CellA.Row > start.CellB.Row { // normalize wall position, make cell A < cell B
			start = utils.WallPosition{CellA: start.CellB, CellB: start.CellA}
		}

		for i := 0; i < g.wallLength; i++ {
			cellAHash := CellHash(Cell{Column: start.CellA.Column + i, Row: start.CellA.Row})
			cellBHash := CellHash(Cell{Column: start.CellB.Column + i, Row: start.CellB.Row})

			if err := _g.AddEdge(cellAHash, cellBHash); err != nil {
				err = fmt.Errorf("Error adding edge between %+v and %+v: %s\n", start.CellA, start.CellB, err)
				return err
			}
		}
	} else {
		if start.CellA.Column > start.CellB.Column { // normalize wall position, make cell A < cell B
			start = utils.WallPosition{CellA: start.CellB, CellB: start.CellA}
		}

		for i := 0; i < g.wallLength; i++ {
			cellAHash := CellHash(Cell{Column: start.CellA.Column, Row: start.CellA.Row + i})
			cellBHash := CellHash(Cell{Column: start.CellB.Column, Row: start.CellB.Row + i})

			if err := _g.AddEdge(cellAHash, cellBHash); err != nil {
				err = fmt.Errorf("Error adding edge between %+v and %+v: %s\n", start.CellA, start.CellB, err)
				return err
			}
		}
	}

	g.Graph = _g

	return nil
}

func (g *Graph) ExistsPath(source, target utils.GridPosition) bool {
	exists := false

	sourceHash := CellHash(Cell{Column: source.Column, Row: source.Row})
	targetHash := CellHash(Cell{Column: target.Column, Row: target.Row})
	err := graph.DFS(g.Graph, sourceHash, func(vertex string) bool {
		if vertex == targetHash {
			exists = true
			return true
		}
		return false
	})

	if err != nil {
		return false
	}
	return exists
}

func (g *Graph) IsLegalMove(source, target utils.GridPosition, opponentPositions []*utils.GridPosition) bool {
	for _, oPos := range opponentPositions {
		if target == *oPos {
			return false
		}

		// el jugador intenta saltear al otro jugador
		if g.IsAdjacent(source, *oPos) && g.IsAdjacent(*oPos, target) {
			return true
		}
	}

	return g.IsAdjacent(source, target)
}

func (g *Graph) GetWalls() []utils.WallPosition {
	// Implementation for getting walls
	return g.Walls
}

func (g *Graph) IsAdjacent(source, target utils.GridPosition) bool {
	_, err := g.Graph.Edge(CellHash(Cell{Column: source.Column, Row: source.Row}), CellHash(Cell{Column: target.Column, Row: target.Row}))
	return !errors.Is(err, graph.ErrEdgeNotFound)
}

func infereWallType(start utils.WallPosition) WallType {
	var wallType WallType
	if start.CellA.Row == start.CellB.Row {
		wallType = Vertical
	} else {
		wallType = Horizontal
	}

	return wallType
}
