import { type LobbyPlayer as LobbyPlayerMessage } from "../../server/messages"

export interface Lobby {
    // The amount of players that can play in this lobby. For now, this will always be 2, but in the future we might want to support more players.
    playerAmount: number,
    
    players: LobbyPlayer[],
    winnerPlayerId?: number
}

type LobbyPlayer = LobbyPlayerMessage 

export const DefaultLobby: Lobby = {
    playerAmount: 2,
    players: [],
    winnerPlayerId: undefined
}