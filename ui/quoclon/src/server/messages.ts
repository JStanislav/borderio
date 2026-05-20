export interface PlayerConfigurationMessage {
    type: string,
    id: number,
    name: string,
    ppid: string
}

export interface LobbyMessage {
    type: string,
    players: LobbyPlayer[]
    winnerPlayerId?: number
}

export interface PlayerJoinedMessage {
    type: string,
    id: number,
    name: string
}

export interface LobbyPlayer {
    id: number,
    name: string,
    ready: boolean
}

export interface GameOverMessage {
    type: string,
    winnerPlayerId: string
}