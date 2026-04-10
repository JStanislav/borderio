
export interface GameState {
  type: string,
  currentTurnPlayerId: number,
  playerOne: {
    id: number,
    name: string,
    position: {
      row: number,
      col: number
    }
    wallsRemaining: number
  },
  playerTwo: {
    id: number,
    name: string,
    position: {
      row: number,
      col: number
    }
    wallsRemaining: number
  }
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

export const getDefaultGameState = (): GameState => {
    return {
      type: "gameState",
      currentTurnPlayerId: 1,
      playerOne: {
        id: 1,
        name: "P1",
        position: {
          row: 0,
          col: 0
        },
        wallsRemaining: 10
      },
      playerTwo: {
        id: 2,
        name: "P2",
        position: {
          row: 8,
          col: 8
        },
        wallsRemaining: 10
      },
      walls: []
    }
}