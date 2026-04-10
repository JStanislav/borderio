package game

import (
	"errors"
	"fmt"
	"time"

	"github.com/JStanislav/quoridor-clone/graph"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
)

type GameState struct {
	StartTime   *time.Time
	CurrentTurn player.PlayerID
	Players     []*player.Player
	WallLength  int

	// Board related
	Board          graph.Board
	FinishLineType FinishLineType //
	Columns        int
	Rows           int
}

type FinishLineType string

const (
	Horizontal            FinishLineType = "horizontal"
	Vertical              FinishLineType = "vertical"
	HorizontalAndVertical FinishLineType = "square"
)

func New(wallLength int, players []*player.Player, columns, rows int, finishLineType FinishLineType) *GameState {
	return &GameState{
		WallLength:     wallLength,
		Players:        players,
		FinishLineType: finishLineType,
		Columns:        columns,
		Rows:           rows,
	}
}

func (g *GameState) StartMatch(movements chan player.Play) {
	boardDimension := 9
	actualBoardDimension := boardDimension + 2

	g.Board = graph.New(2, graph.ExtraRows)
	g.Board.GenerateBoard(boardDimension, actualBoardDimension)

	g.StartTime = new(time.Time)
	*g.StartTime = time.Now()

	for _, p := range g.Players {
		playersButNotCurrent := []*player.Player{}
		for _, other := range g.Players {
			if other.ID != p.ID {
				playersButNotCurrent = append(playersButNotCurrent, other)
			}
		}
		p.OnPlayerPlay = func(playerID player.PlayerID, play player.Play) error {

			if p.ID != g.CurrentTurn {
				fmt.Printf("Player %d attempted to place wall out of turn\n", p.ID)
				return errors.New("not your turn")
			}

			playersButNotCurrentPositions := player.GetPlayersPositions(playersButNotCurrent)

			switch play.PlayType {
			case player.PlayerMove:
				fmt.Printf("Moving P%d [R%d-C%d]->[R%d-C%d]\n", p.ID, p.Position.Row, p.Position.Column, play.Position.Row, play.Position.Column)

				if g.OutOfBounds(play, boardDimension, actualBoardDimension) && !p.IsFinishLine(play) {
					return errors.New("move out of bounds")
				}

				if g.Board.IsLegalMove(*p.Position, *play.Position, playersButNotCurrentPositions) {
					p.Position = play.Position
					g.CurrentTurn = playersButNotCurrent[0].ID // WRONG! IMPLEMENT INFINITE STATE MACHINE
					movements <- play
					fmt.Println("Moved")

					if p.IsWinner() {
						fmt.Printf("Player %d wins!\n", p.ID)
					}

					return nil
				} else {
					return errors.New("illegal move")
				}

			case player.WallPlacement:
				if p.WallsRemaining <= 0 {
					return errors.New("player has no more walls")
				}
				wallPosition := utils.WallPosition{CellA: play.WallPlaced.CellA, CellB: play.WallPlaced.CellB}

				fmt.Printf("Placing wall p%d [R%d-C%d]||[R%d-C%d]\n", p.ID, play.WallPlaced.CellA.Row, play.WallPlaced.CellA.Column, play.WallPlaced.CellB.Row, play.WallPlaced.CellB.Column)

				if g.OutOfBounds(play, boardDimension, actualBoardDimension) {
					return errors.New("wall out of bounds")
				}

				err := g.Board.AddWall(graph.Undefined, wallPosition)
				if err != nil {
					return errors.New("illegal wall placement")
				}

				finishLinesFound := 0
				for _, p := range g.Players {
					existsPathToFinishLine := false
					if p.FinishLine.Type == utils.HorizontalLine {
						for i := range boardDimension {
							winCell := utils.GridPosition{Row: p.FinishLine.Index, Column: i}
							if g.Board.ExistsPath(*p.Position, winCell) {
								existsPathToFinishLine = true
								break
							}
						}
					} else {
						for i := 1; i < actualBoardDimension-1; i++ {
							winCell := utils.GridPosition{Row: i, Column: p.FinishLine.Index}
							if g.Board.ExistsPath(*p.Position, winCell) {
								existsPathToFinishLine = true
								break
							}
						}
					}
					if existsPathToFinishLine {
						finishLinesFound++
					}
				}

				if finishLinesFound != len(g.Players) {
					g.Board.RemoveWall(graph.Undefined, wallPosition)
					return errors.New("illegal wall placement, no path to finish line")
				}

				g.CurrentTurn = playersButNotCurrent[0].ID // WRONG! IMPLEMENT INFINITE STATE MACHINE
				movements <- play
				p.WallsRemaining -= 1
				fmt.Println("Placed wall")

				return nil
			}

			return errors.New("unexpected error ")
		}
	}

}

func RowOutOfBounds(p utils.GridPosition, row int) bool {
	return p.Row < 1 || p.Row >= row-1
}

func ColumnOutOfBounds(p utils.GridPosition, column int) bool {
	return p.Column < 0 || p.Column >= column
}

func (g *GameState) OutOfBounds(p player.Play, columns, rows int) bool {
	switch p.PlayType {
	case player.PlayerMove:
		return RowOutOfBounds(*p.Position, rows) || ColumnOutOfBounds(*p.Position, columns)
	case player.WallPlacement:
		if p.WallPlaced.Orientation() == utils.VerticalLine && RowOutOfBounds(p.WallPlaced.CellA, rows-(g.WallLength-1)) {
			return true
		}
		if p.WallPlaced.Orientation() == utils.HorizontalLine && ColumnOutOfBounds(p.WallPlaced.CellA, columns-(g.WallLength-1)) {
			return true
		}
		return false
	default:
		return true
	}
}
