package game

import (
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
func (g *GameState) StartMatch(playerOne, playerTwo *player.Player, p1Moves chan utils.GridPosition, p2Moves chan utils.GridPosition, movements chan<- PlayerMovement) {
	g.Board = graph.New()
	p1StartPosition := utils.GridPosition{Column: 4, Row: 0}
	p2StartPosition := utils.GridPosition{Column: 4, Row: 8}

	g.Board.GenerateBoard(9, 9, p1StartPosition, p2StartPosition)

	g.StartTime = new(time.Time)
	*g.StartTime = time.Now()
	g.CurrentTurn = playerOne.ID

	for {
		select {
		case move := <-p1Moves:
			fmt.Printf("Moving p1 [R%d-C%d]->[R%d-C%d]\n", playerOne.Position.Row, playerOne.Position.Column, move.Row, move.Column)
			if g.Board.IsLegalMove(playerOne.Position, move, playerTwo.Position) && g.CurrentTurn == playerOne.ID {
				playerOne.Position = move
				g.CurrentTurn = playerTwo.ID
				movements <- PlayerMovement{PlayerID: playerOne.ID, Move: &playerOne.Position}
			} else {
				movements <- PlayerMovement{PlayerID: playerOne.ID, Move: nil}
			}

		case move := <-p2Moves:
			fmt.Printf("Moving p2 [R%d-C%d]->[R%d-C%d]\n", playerTwo.Position.Row, playerTwo.Position.Column, move.Row, move.Column)
			if g.Board.IsLegalMove(playerTwo.Position, move, playerOne.Position) && g.CurrentTurn == playerTwo.ID {
				playerTwo.Position = move
				g.CurrentTurn = playerOne.ID
				movements <- PlayerMovement{PlayerID: playerTwo.ID, Move: &playerTwo.Position}
			} else {
				movements <- PlayerMovement{PlayerID: playerTwo.ID, Move: nil}
			}
		}
	}
}

type PlayerMovement struct {
	PlayerID player.PlayerID
	Move     *utils.GridPosition
}
