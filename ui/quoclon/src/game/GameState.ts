
export interface GameState {
  type: string,
  playerOne: {
    id: number,
    name: string,
    position: {
      row: number,
      col: number
    }
  },
  playerTwo: {
    id: number,
    name: string,
    position: {
      row: number,
      col: number
    }
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
      playerOne: {
        id: 1,
        name: "P1",
        position: {
          row: 0,
          col: 0
        }
      },
      playerTwo: {
        id: 2,
        name: "P2",
        position: {
          row: 8,
          col: 8
        }
      },
      walls: []
    }
}