import type { GameState } from "./App"
import { Board } from "./components/Board/NewBoard"
import { WallPicker } from "./components/Board/WallPicker"
import "./gameframe.css"
import { send } from "./server-conn"

const translateGridPositionToServer = (row: number, col: number) => {
    return {
        s_row: Math.floor(row / 2),
        s_col: Math.floor(col / 2)
    }
}

const translateGridPositionToClient = (s_row: number, s_col: number) => {
    return {
        row: s_row * 2 + 1,
        col: s_col * 2 + 1
    }
}

export const GameFrame = ({ gameState }: { gameState: GameState }) => {
    let p1Position = translateGridPositionToClient(gameState.playerOne.position.row, gameState.playerOne.position.col);
    let p2Position = translateGridPositionToClient(gameState.playerTwo.position.row, gameState.playerTwo.position.col);
    const players = [
        {
            id: gameState.playerOne.id,
            name: gameState.playerOne.name,
            position: {
                row: p1Position.row,
                col: p1Position.col
            }
        },
        {
            id: gameState.playerTwo.id,
            name: gameState.playerTwo.name,
            position: {
                row: p2Position.row,
                col: p2Position.col
            }
        }
    ]

    const requestPlayerMove = (playerId: number, row: number, col: number) => {
        const { s_row, s_col } = translateGridPositionToServer(row, col);
        const type = "playerMove";
        const target = { row: s_row,col: s_col };
        send(type, { playerId, target });
    }

    return (
        <div className="game-frame">
            <WallPicker walls={9} position="top"/>
            <Board players={players} requestPlayerMove={requestPlayerMove}/>
            <WallPicker walls={9} position="bottom"/>
        </div>
    )
}