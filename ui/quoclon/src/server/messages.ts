export interface PlayerConfigurationMessage {
    type: string,
    id: number,
    name: string,
    ppid: string
}

export interface LobbyMessage {
    type: string,
    players: LobbyPlayer[]
}

export interface LobbyPlayer {
    id: number,
    name: string,
    ready: boolean
}