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
}

func New() *GameState {
	return &GameState{}
}

func (g *GameState) StartMatch(playerOne, playerTwo *player.Player, movements chan player.Play) {
	g.Board = graph.New(2)
	p1StartPosition := utils.GridPosition{Column: 4, Row: 0}
	p2StartPosition := utils.GridPosition{Column: 4, Row: 8}

	g.Board.GenerateBoard(9, 9, p1StartPosition, p2StartPosition)

	g.StartTime = new(time.Time)
	*g.StartTime = time.Now()
	g.CurrentTurn = playerOne.ID

	playerOne.Position = p1StartPosition
	playerTwo.Position = p2StartPosition

	playerOne.OnPlayerPlay = func(playerID player.PlayerID, play player.Play) error {
		switch play.PlayType {
		case player.PlayerMove:
			fmt.Printf("Moving p1 [R%d-C%d]->[R%d-C%d]\n", playerOne.Position.Row, playerOne.Position.Column, play.Position.Row, play.Position.Column)
			if g.Board.IsLegalMove(playerOne.Position, *play.Position, playerTwo.Position) && g.CurrentTurn == playerOne.ID {
				playerOne.Position = *play.Position
				g.CurrentTurn = playerTwo.ID
				movements <- play
				fmt.Println("Moved")
				return nil
			} else {
				return errors.New("illegal move")
			}
		case player.WallPlacement:
			fmt.Printf("Placing wall p1 [R%d-C%d]||[R%d-C%d]\n", play.WallPlaced.CellA.Row, play.WallPlaced.CellA.Column, play.WallPlaced.CellB.Row, play.WallPlaced.CellB.Column)
			if g.Board.AddWall(graph.Undefined, utils.WallPosition{CellA: play.WallPlaced.CellA, CellB: play.WallPlaced.CellB}) == nil && g.CurrentTurn == playerOne.ID {
				movements <- play
				fmt.Println("Placed wall")
				return nil
			} else {
				return errors.New("illegal wall placement")
			}
		}

		return errors.New("unexpected error ")
	}

	playerTwo.OnPlayerPlay = func(playerID player.PlayerID, play player.Play) error {
		switch play.PlayType {
		case player.PlayerMove:
			fmt.Printf("Moving p2 [R%d-C%d]->[R%d-C%d]\n", playerTwo.Position.Row, playerTwo.Position.Column, play.Position.Row, play.Position.Column)
			if g.Board.IsLegalMove(playerTwo.Position, *play.Position, playerOne.Position) && g.CurrentTurn == playerTwo.ID {
				playerTwo.Position = *play.Position
				g.CurrentTurn = playerOne.ID
				movements <- play
				fmt.Println("Moved")
				return nil
			} else {
				return errors.New("illegal move")
			}
		case player.WallPlacement:
			fmt.Printf("Placing wall p2 [R%d-C%d]||[R%d-C%d]\n", play.WallPlaced.CellA.Row, play.WallPlaced.CellA.Column, play.WallPlaced.CellB.Row, play.WallPlaced.CellB.Column)
			if g.Board.AddWall(graph.Undefined, utils.WallPosition{CellA: play.WallPlaced.CellA, CellB: play.WallPlaced.CellB}) == nil && g.CurrentTurn == playerTwo.ID {
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
