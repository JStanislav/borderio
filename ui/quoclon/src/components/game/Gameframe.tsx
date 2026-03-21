import { WallPicker } from "../board/WallPicker"
import "./gameframe.css"
import { type GameState } from "../../game/GameState"
import { requestPlayerMove, requestWallPlacement } from "../../server/server"
import { translateGridPositionToClient, translateWallsToClient } from "../../server/utils"
import { Board } from "../board/Board"


export const GameFrame = ({ gameState }: { gameState: GameState }) => {
    const p1Position = translateGridPositionToClient(gameState.playerOne.position.row, gameState.playerOne.position.col);
    const p2Position = translateGridPositionToClient(gameState.playerTwo.position.row, gameState.playerTwo.position.col);
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

    const activeWalls = translateWallsToClient(gameState.walls || []);

    return (
        <div className="game-frame">
            <WallPicker walls={9} position="top"/>
            <Board players={players} requestPlayerMove={requestPlayerMove} requestWallPlacement={requestWallPlacement} activeWalls={activeWalls}/>
            <WallPicker walls={9} position="bottom"/>
        </div>
    )
}