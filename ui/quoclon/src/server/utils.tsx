


export const translateGridPositionToServer = (row: number, col: number) => {
    return {
        s_row: Math.floor(row / 2),
        s_col: Math.floor(col / 2)
    }
}

export const translateGridPositionToClient = (s_row: number, s_col: number) => {
    return {
        row: s_row * 2 + 1,
        col: s_col * 2 + 1
    }
}

interface s_Position {
    row: number,
    column: number
}

interface s_WallTarget {
    cellA: s_Position,
    cellB: s_Position,
}

export const translateWallGridPositionToServer = (row: number, col: number, orientation: "horizontal" | "vertical"): s_WallTarget  => {
    const positions = translateGridPositionToServer(row, col);
    if (orientation === "horizontal") {
        return {
            cellA: {
                row: positions.s_row,
                column: positions.s_col
            },
            cellB: {
                row: positions.s_row -1,
                column: positions.s_col
            }
        }
    } else {
        return {
            cellA: {
                row: positions.s_row,
                column: positions.s_col
            },
            cellB: {
                row: positions.s_row,
                column: positions.s_col -1
            }
        }
    }
}

export const translateWallGridPositionToClient = (s_cellA: s_Position, s_cellB: s_Position): { row: number, col: number, orientation: "horizontal" | "vertical" } => {
    const cellAClientPos = translateGridPositionToClient(s_cellA.row, s_cellA.column);
    const cellBClientPos = translateGridPositionToClient(s_cellB.row, s_cellB.column);
    const orientation = cellAClientPos.row === cellBClientPos.row ? "vertical" : "horizontal";
    return {
        row: orientation === "vertical" ? cellAClientPos.row : cellAClientPos.row + 1,
        col: orientation === "vertical" ? cellAClientPos.col + 1 : cellAClientPos.col ,
        orientation
    }
}

export const translateWallsToClient = (walls: Array<{ cellA: s_Position, cellB: s_Position }>) => {
    return walls.map(wall => translateWallGridPositionToClient(wall.cellA, wall.cellB))
}