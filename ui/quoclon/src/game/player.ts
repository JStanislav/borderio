export interface Player {
    id: number,
    name: string,
    ppid: string,

    inGame?: boolean,
}

export const DefaultPlayer: Player = {
    id: -1,
    name: "",
    ppid: "",
    inGame: false,
}