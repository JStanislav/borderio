
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

/* 
  TODO: fix this semantics.
  The function is used to know whether a game started and not to know if all players are ready. The name is misleading.
  But also the gameState.playerOne.ready and gameState.playerTwo.ready are wrong, since their are set to true when the game starts, 
  and not when the players are ready. So this function is not really useful, since it will always return true when the game starts.
*/
export const allPlayersReady = (gameState: GameState): boolean => {
  return gameState.playerOne.ready && gameState.playerTwo.ready;
}

export const getPlayerById = (gameState: GameState, playerId: number): Player | undefined => {
  if (gameState.playerOne.id === playerId) {
    return gameState.playerOne;
  }
  if (gameState.playerTwo.id === playerId) {
    return gameState.playerTwo;
  }
  return undefined;
}