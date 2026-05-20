import { type LobbyPlayer as LobbyPlayerMessage } from "../../server/messages"

export interface Lobby {
    players: LobbyPlayer[],
    winnerPlayerId?: number
}

type LobbyPlayer = LobbyPlayerMessage 

export const DefaultLobby: Lobby = {
    players: [],
    winnerPlayerId: undefined
}