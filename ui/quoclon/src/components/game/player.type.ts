export type Players = Player[]

export type Player = {
    id: number,
    name: string,
    position: {
        row: number,
        col: number
    }
}
