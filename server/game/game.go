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
	Board       graph.Board
	StartTime   *time.Time
	CurrentTurn player.PlayerID
	Players     []*player.Player
}

func New() *GameState {
	return &GameState{}
}

func (g *GameState) StartMatch(playerOne, playerTwo *player.Player, movements chan player.Play) {
	boardDimension := 9
	actualBoardDimension := boardDimension + 2

	g.Board = graph.New(2)
	p1StartPosition := utils.GridPosition{Column: 4, Row: 0}
	p2StartPosition := utils.GridPosition{Column: 4, Row: actualBoardDimension - 1}

	p1StartLine := utils.Line{Type: utils.HorizontalLine, Index: p1StartPosition.Row}
	p2StartLine := utils.Line{Type: utils.HorizontalLine, Index: p2StartPosition.Row}
	p1FinishLine := utils.Line{Type: utils.HorizontalLine, Index: p2StartPosition.Row}
	p2FinishLine := utils.Line{Type: utils.HorizontalLine, Index: p1StartPosition.Row}

	g.Board.GenerateBoard(boardDimension, actualBoardDimension, p1StartPosition, p2StartPosition)

	g.StartTime = new(time.Time)
	*g.StartTime = time.Now()
	g.CurrentTurn = playerOne.ID

	playerOne.Position = &p1StartPosition
	playerOne.StartLine = p1StartLine
	playerOne.FinishLine = p1FinishLine
	playerTwo.Position = &p2StartPosition
	playerTwo.StartLine = p2StartLine
	playerTwo.FinishLine = p2FinishLine

	g.Players = []*player.Player{playerOne, playerTwo}

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

				if play.OutOfBounds(boardDimension, actualBoardDimension) && !p.IsFinishLine(play) {
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
				fmt.Printf("Placing wall p%d [R%d-C%d]||[R%d-C%d]\n", p.ID, play.WallPlaced.CellA.Row, play.WallPlaced.CellA.Column, play.WallPlaced.CellB.Row, play.WallPlaced.CellB.Column)

				if g.Board.AddWall(graph.Undefined, utils.WallPosition{CellA: play.WallPlaced.CellA, CellB: play.WallPlaced.CellB}) == nil {
					movements <- play
					fmt.Println("Placed wall")
					return nil
				} else {
					return errors.New("illegal wall placement")
				}
			}

			return errors.New("unexpected error ")
		}
	}

}
