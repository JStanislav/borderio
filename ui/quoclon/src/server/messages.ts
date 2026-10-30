export interface IMessage<T> {
    type: string
    payload: T
}

export interface PlayerConfigurationMessage {
    id: number,
    name: string,
    ppid: string
}

export interface LobbyMessage {
    id: string,
    players: LobbyPlayer[]
    gameStarted: boolean
    winnerPlayerId?: number
    draw?: boolean
}

export interface PlayerJoinedMessage {
    id: number,
    name: string
}

export interface LobbyPlayer {
    id: number,
    name: string,
    ready: boolean
    host: boolean
    connected: boolean
    spectator: boolean
}

export interface GameOverMessage {
    winnerPlayerId: string
}

export interface PlayerDisconnectedMessage {
    id: number,
    name: string,
    gracePeriod: number // in seconds
}

export interface WillTimeoutMessage {
    span: string,
}
