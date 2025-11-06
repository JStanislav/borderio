type RowType = "Full" | "WallOnly"

export interface Row {
    id: number,
    type: RowType
    cells: Array<Cell>
}

type CellType = "Square" | "Wall" | "Blank"

export interface Cell {
    id: string
    type: CellType
}

const BOARD_DIM = 3
const TOTAL_BOARD_DIM = BOARD_DIM * 2 + 1

export const board: Array<Row> = Array.from({length: BOARD_DIM * 2 + 1}, (_, row) => {
    const cells: Array<Cell> = Array.from({length: BOARD_DIM * 2 + 1}, (_, col) => {
        let type: CellType
        if ((row + col) % 2 === 0) { // sum is even
            if ((row % 2 === 0) && (col % 2 === 0)) { // even + even
                type = "Blank"
            } else {
                type = "Square"
            } 
        } else { // sum is odd
            type = "Wall"
        }
        return {
            id: `cell-${(TOTAL_BOARD_DIM*row) + col}`,
            type: type
        }
    })

    return {
        id: row,
        type: row % 2 === 0 ? "WallOnly" : "Full",
        cells
    }
})



