

interface CellProps {
    topWalled: boolean,
    bottomWalled: boolean,
    leftWalled: boolean,
    rightWalled: boolean,
    type: string
}


const generateBoard = (rows: number, columns: number) => {

    let board: CellProps[][] = []

    for (let i = 0; i <= rows; i++) {
        board[i] = []
        for (let j = 0; j <= columns; j++) {
            board[i][j] = {
                topWalled: false,
                bottomWalled: false,
                leftWalled: false,
                rightWalled: false,
                type: "box"
            }
        }
    }

    board[4][3].bottomWalled = true
    board[5][3].topWalled = true

    board[2][3].leftWalled = true
    board[2][2].rightWalled = true
    board[1][3].leftWalled = true
    board[1][2].rightWalled = true

    return board
}


export const getBoard = (): CellProps[][] => {
    return generateBoard(9, 9)
}