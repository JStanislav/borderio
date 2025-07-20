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

/*
We receive player moves from the channels and send them to the movements channel
If a player makes an illegal move, we send a nil to the movements channel
*/
func (g *GameState) StartMatch(playerOne, playerTwo *player.Player, movements chan<- player.Play) {
	g.Board = graph.New()
	p1StartPosition := utils.GridPosition{Column: 4, Row: 0}
	p2StartPosition := utils.GridPosition{Column: 4, Row: 8}

	g.Board.GenerateBoard(9, 9, p1StartPosition, p2StartPosition)

	g.StartTime = new(time.Time)
	*g.StartTime = time.Now()
	g.CurrentTurn = playerOne.ID

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
			if g.Board.AddWall(playerOne.Position.Column, playerOne.Position.Row, play.Position.Column, play.Position.Row) == nil && g.CurrentTurn == playerOne.ID {
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
			if g.Board.AddWall(playerTwo.Position.Column, playerTwo.Position.Row, play.Position.Column, play.Position.Row) == nil && g.CurrentTurn == playerTwo.ID {
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
