
export interface GameState {
  type: string,
  currentTurnPlayerId: number,
  playerOne: Player,
  playerTwo: Player,
  walls: Array<{
    cellA: {
      row: number,
      column: number
    },
    cellB: {
      row: number,
      column: number
    },
  }>
}

interface Player {
  id: number,
  name: string,
  position: {
    row: number,
    col: number
  }
  wallsRemaining: number,
  ready: boolean
}

export const getDefaultGameState = (): GameState => {
    return {
      type: "gameState",
      currentTurnPlayerId: 1,
      playerOne: {
        id: -1,
        name: "P1",
        position: {
          row: 0,
          col: 0
        },
        wallsRemaining: 10,
        ready: false
      },
      playerTwo: {
        id: -1,
        name: "P2",
        position: {
          row: 8,
          col: 8
        },
        wallsRemaining: 10,
        ready: false
      },
      walls: []
    }
}

export const allPlayersReady = (gameState: GameState): boolean => {
  return gameState.playerOne.ready && gameState.playerTwo.ready;
}