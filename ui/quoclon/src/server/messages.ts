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
    players: LobbyPlayer[]
    winnerPlayerId?: number
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
}

export interface GameOverMessage {
    winnerPlayerId: string
}
