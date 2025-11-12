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
    fillType: "vertical-start" | "vertical-middle" | "vertical-end" | "horizontal-start" | "horizontal-middle" | "horizontal-end" | false;
}

const BOARD_DIM = 9
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
            type: type,
            fillType: false
        }
    })

    return {
        id: row,
        type: row % 2 === 0 ? "WallOnly" : "Full",
        cells
    }
})

export const setWall = (board: Array<Row>, row: number, col: number, wallLength: number = 2) => {
    if (wallLength === 0) return
    if ((row % 2) === 0 && (col % 2) === 0) return
    if ((row % 2) === 1 && (col % 2) === 1) return

    wallLength = wallLength * 2 - 1

    if ((row % 2) === 0) { // horizontal
        // TODO: add validation wall length + col < board limit
        board[row].cells[col].fillType = "horizontal-start"
        for (let i = 1; i < wallLength - 1 ; i++) {
            board[row].cells[col + i].fillType = "horizontal-middle"
        }
        board[row].cells[col + wallLength - 1].fillType = "horizontal-end"

    } else { // vertical
        // TODO: add validation wall length + row < board limit
        board[row].cells[col].fillType = "vertical-start"
        for (let i = 1; i < wallLength-1; i++) {
            board[row + i].cells[col].fillType = "vertical-middle"
        }
        board[row + wallLength - 1].cells[col].fillType = "vertical-end"
    }
}

// dummy test
setWall(board, 4, 3, 2);
setWall(board, 4, 7, 2);
setWall(board, 4, 15, 2);
setWall(board, 4, 11, 2);

setWall(board, 9, 4, 2)
setWall(board, 5, 4, 2)
setWall(board, 1, 4, 2)
setWall(board, 13, 4, 2)
setWall(board, 3, 6, 2)


