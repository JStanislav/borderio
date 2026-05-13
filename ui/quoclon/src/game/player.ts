export interface Player {
    id: number,
    name: string,
    ppid: string,
    ready: boolean
}

export const DefaultPlayer: Player = {
    id: -1,
    name: "",
    ppid: "",
    ready: false
}