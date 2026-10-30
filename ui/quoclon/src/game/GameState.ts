export interface GameState {
  type: string,
  currentTurnPlayerId: number,
  players: Map<number, Player>
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

export interface Player {
  id: number,
  name: string,
  position: {
    row: number,
    col: number
  }
  wallsRemaining: number,
  ready: boolean
}

export const GetDefaultPlayer = (id: number, name: string): Player => {
  return {
    id: id,
    name: name,
    position: {
      row: 0,
      col: 0
    },
    wallsRemaining: 0,
    ready: false
  }
}

export const getDefaultGameState = (): GameState => {
    return {
      type: "gameState",
      currentTurnPlayerId: 1,
      players: new Map(),
      walls: []
    }
}