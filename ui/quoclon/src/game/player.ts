export interface Player {
    id: number,
    name: string,
    ppid: string,
}

export const DefaultPlayer: Player = {
    id: -1,
    name: "",
    ppid: "",
}