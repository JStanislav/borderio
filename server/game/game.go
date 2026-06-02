package game

import (
	"container/ring"
	"errors"
	"fmt"
	"time"

	"github.com/JStanislav/quoridor-clone/graph"
	"github.com/JStanislav/quoridor-clone/player"
	"github.com/JStanislav/quoridor-clone/utils"
)

type GameState struct {
	PlayerCount int
	StartTime   *time.Time
	Players     *[]*player.Player
	WallLength  int

	// Board related
	Board          graph.Board
	FinishLineType FinishLineType //
	Columns        int
	Rows           int

	Turner *ring.Ring
}

type FinishLineType string

const (
	Horizontal            FinishLineType = "horizontal"
	Vertical              FinishLineType = "vertical"
	HorizontalAndVertical FinishLineType = "square"
)

func New(wallLength int, players int, columns, rows int, finishLineType FinishLineType) *GameState {
	board := graph.New(2, graph.ExtraRows)
	return &GameState{
		PlayerCount:    players,
		WallLength:     wallLength,
		FinishLineType: finishLineType,
		Columns:        columns,
		Rows:           rows,
		Players:        &[]*player.Player{},
		StartTime:      new(time.Time),
		Turner:         ring.New(players),
		Board:          board,
	}
}

func (g *GameState) StartMatchWithMovementsChannel() chan player.Play {
	movements := make(chan player.Play)

	go func() {
		for mov := range movements {
			fmt.Printf("Received movement: %+v\n", mov)
		}
	}()

	g.StartMatch(movements)

	return movements
}

func (g *GameState) StartMatch(movements chan player.Play) {
	boardDimension := 9
	actualBoardDimension := boardDimension + 2

	g.Board.GenerateBoard(boardDimension, actualBoardDimension)

	*g.StartTime = time.Now()

	for _, p := range *g.Players {
		g.Turner.Value = p
		g.Turner = g.Turner.Next()
	}

	for _, p := range *g.Players {
		playersButNotCurrent := g.GetPlayersExcept(p.ID)

		p.OnPlayerPlay = func(playerID player.PlayerID, play player.Play) error {

			if p.ID != g.GetCurrentTurnPlayer().ID {
				fmt.Printf("Player %d attempted to play out of their turn\n", p.ID)
				return errors.New("not your turn")
			}

			playersButNotCurrentPositions := player.GetPlayersPositions(*playersButNotCurrent)

			if g.OutOfBounds(play, boardDimension, actualBoardDimension) && !p.IsMovingToFinishLine(play) {
				return errors.New("play out of bounds")
			}

			switch play.PlayType {
			case player.PlayerMove:

				fmt.Printf("Moving P%d [R%d-C%d]->[R%d-C%d]\n", p.ID, p.Position.Row, p.Position.Column, play.Position.Row, play.Position.Column)

				if g.Board.IsLegalMove(*p.Position, *play.Position, playersButNotCurrentPositions) {
					p.Position = play.Position
					g.Turner = g.Turner.Next()
					movements <- play
					fmt.Println("Moved")

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

				err := g.Board.AddWall(graph.Undefined, wallPosition)
				if err != nil {
					return errors.New("illegal wall placement")
				}

				if !g.PlayersCanReachFinishLine(boardDimension, actualBoardDimension) {
					g.Board.RemoveWall(graph.Undefined, wallPosition)
					return errors.New("illegal wall placement, prevents players from reaching finish line")
				}

				g.Turner = g.Turner.Next()
				movements <- play
				p.WallsRemaining -= 1
				fmt.Println("Placed wall")

				return nil
			}

			return errors.New("unexpected error ")
		}
	}

}

// Checks if all players can reach their finish line. Should be called after every wall placement to ensure the game is still winnable.
// Works for any number of players.
// Works for both horizontal and vertical finish lines
func (g *GameState) PlayersCanReachFinishLine(columns, rows int) bool {
	finishLinesFound := 0
	for _, p := range *g.Players {
		existsPathToFinishLine := false
		if p.FinishLine.Type == utils.HorizontalLine {
			for i := range columns {
				winCell := utils.GridPosition{Row: p.FinishLine.Index, Column: i}
				if g.Board.ExistsPath(*p.Position, winCell) {
					existsPathToFinishLine = true
					break
				}
			}
		} else {
			for i := 1; i < rows-1; i++ {
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

	return finishLinesFound == len(*g.Players)
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

func (g *GameState) GetCurrentTurnPlayer() *player.Player {
	if g.Turner.Value != nil {
		return g.Turner.Value.(*player.Player)
	} else {
		return nil
	}
}

func (g *GameState) AllPlayersReady() bool {
	for _, p := range *g.Players {
		if !p.Ready {
			return false
		}
	}
	return true
}

func (g *GameState) AddPlayer(p *player.Player) error {
	if len(*g.Players) < g.PlayerCount {
		*g.Players = append(*g.Players, p)

		if len(*g.Players) == 1 {
			p.Host = true
		}
		return nil
	}

	return errors.New("game is full")
}

func (g *GameState) RemovePlayer(pid player.PlayerID) error {
	for i, player := range *g.Players {
		if player.ID == pid {
			var wasHost bool
			if player.Host {
				wasHost = true
			}
			*g.Players = append((*g.Players)[:i], (*g.Players)[i+1:]...)
			if wasHost && len(*g.Players) > 0 {
				(*g.Players)[0].Host = true
			}
			return nil
		}
	}
	return errors.New("player not found")
}

func (g *GameState) GetPlayerPPID(ppid string) *player.Player {
	for _, p := range *g.Players {
		if p.PrivatePlayerID == ppid {
			return p
		}
	}
	return nil
}

func (g *GameState) GetGameStats() GameStats {
	var winnerId int
	if (*g.Players)[0].IsWinner() {
		winnerId = int((*g.Players)[0].ID)
	} else if (*g.Players)[1].IsWinner() {
		winnerId = int((*g.Players)[1].ID)
	}

	return GameStats{
		PlayerWinnerId:      winnerId,
		TotalMoves:          0,
		TotalWallPlacements: 0,
		TotalWallsRemaining: 0,
		StartTime:           &time.Time{},
		EndTime:             &time.Time{},
		Players:             &[]*player.Player{},
		Points:              0,
		Steps:               []player.Play{},
	}
}

func (g *GameState) GetUnusedPlayerID() int {
	id := 1
	for _, player := range *g.Players {
		if int(player.ID) != id {
			break
		}
		id++
	}
	return id
}

func (g *GameState) GetPlayersExcept(pid player.PlayerID) *[]*player.Player {
	players := make([]*player.Player, 0)
	for _, p := range *g.Players {
		if p.ID != pid {
			players = append(players, p)
		}
	}
	return &players
}
