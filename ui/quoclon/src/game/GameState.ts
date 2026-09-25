
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

export const getDefaultGameState = (): GameState => {
    return {
      type: "gameState",
      currentTurnPlayerId: 1,
      players: new Map(),
      walls: []
    }
}

/* 
  TODO: fix this semantics.
  The function is used to know whether a game started and not to know if all players are ready. The name is misleading.
  But also the gameState.playerOne.ready and gameState.playerTwo.ready are wrong, since their are set to true when the game starts, 
  and not when the players are ready. So this function is not really useful, since it will always return true when the game starts.
*/
export const allPlayersReady = (gameState: GameState): boolean => {
  if (gameState.players.size === 0) {
    return false
  }

  let allReady = true

  gameState.players.forEach((value) => {
    if(!value) return
    if (!value.ready) {
      allReady = false
    }
  })

  return allReady
}