export type Players = Player[]

export type Player = {
    id: number,
    name: string,
    position: {
        row: number,
        col: number
    }
}

export function getGameOverText(winnerPlayerName: string, winner: boolean, draw: boolean = false): string {
    if (draw) {
        return `The game ended in a draw`
    }
    if (winner) {
        return `Congratulations! You won!`
    } else {
        return `The winner is ${winnerPlayerName}`
    }
}