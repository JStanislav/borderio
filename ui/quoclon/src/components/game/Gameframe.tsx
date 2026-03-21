import { Board } from "../Board/NewBoard"
import { WallPicker } from "../Board/WallPicker"
import "./gameframe.css"
import { send } from "../../server/server-conn"
import type { GameState } from "../../game/GameState"

const translateGridPositionToServer = (row: number, col: number) => {
    return {
        s_row: Math.floor(row / 2),
        s_col: Math.floor(col / 2)
    }
}

const translateGridPositionToClient = (s_row: number, s_col: number) => {
    console.log("translating grid position to client with s_row", s_row, "s_col", s_col);
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

const translateWallGridPositionToServer = (row: number, col: number, orientation: "horizontal" | "vertical"): s_WallTarget  => {
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

const translateWallGridPositionToClient = (s_cellA: s_Position, s_cellB: s_Position): { row: number, col: number, orientation: "horizontal" | "vertical" } => {
    console.log("translating wall position to client with s_cellA", s_cellA, "s_cellB", s_cellB);
    const cellAClientPos = translateGridPositionToClient(s_cellA.row, s_cellA.column);
    const cellBClientPos = translateGridPositionToClient(s_cellB.row, s_cellB.column);
    console.debug("cellA client pos", cellAClientPos, "cellB client pos", cellBClientPos);
    const orientation = cellAClientPos.row === cellBClientPos.row ? "vertical" : "horizontal";
    return {
        row: orientation === "vertical" ? cellAClientPos.row : cellAClientPos.row + 1,
        col: orientation === "vertical" ? cellAClientPos.col + 1 : cellAClientPos.col ,
        orientation
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

    const activeWalls = gameState.walls?.map(wall => translateWallGridPositionToClient(wall.cellA, wall.cellB))

    const requestPlayerMove = (playerId: number, row: number, col: number) => {
        const { s_row, s_col } = translateGridPositionToServer(row, col);
        const type = "playerMove";
        const target = { row: s_row,col: s_col };
        send(type, { playerId, target });
    }

    const requestWallPlacement = (playerId: number, row: number, col: number, orientation: "horizontal" | "vertical") => {
        let wallPositions = translateWallGridPositionToServer(row, col, orientation);
        const type = "wallPlacement";
        if (orientation === "horizontal") {
            if (wallPositions.cellA.row > wallPositions.cellB.row) {
                wallPositions = { cellA: wallPositions.cellB, cellB: wallPositions.cellA };
            }
        }
        const wallTarget = { cellA: { row: wallPositions.cellA.row, col: wallPositions.cellA.column }, cellB: { row: wallPositions.cellB.row, col: wallPositions.cellB.column }, orientation };
        console.log("requesting wall placement with target", wallTarget);
        send(type, { playerId, wallTarget: wallTarget });
    }

    return (
        <div className="game-frame">
            <WallPicker walls={9} position="top"/>
            <Board players={players} requestPlayerMove={requestPlayerMove} requestWallPlacement={requestWallPlacement} activeWalls={activeWalls}/>
            <WallPicker walls={9} position="bottom"/>
        </div>
    )
}