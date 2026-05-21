import { WallPicker } from "../board/WallPicker"
import "./gameframe.css"
import { getPlayerById, type GameState } from "../../game/GameState"
import { requestPlayerMove, requestWallPlacement } from "../../server/server"
import { translateGridPositionToClient, translateWallsToClient } from "../../server/utils"
import { Board } from "../board/Board"
import { GameOverDialog } from "./GameOverDialog.tsx"
import { useContext } from "react"
import { LobbyContext } from "../../App.tsx"


export const GameFrame = ({ gameState }: { gameState: GameState }) => {
    const lobbyContext = useContext(LobbyContext);

    const p1Position = translateGridPositionToClient(gameState.playerOne.position.row, gameState.playerOne.position.col);
    const p2Position = translateGridPositionToClient(gameState.playerTwo.position.row, gameState.playerTwo.position.col);
    const players = [
        {
            id: gameState.playerOne.id,
            name: gameState.playerOne.name,
            position: p1Position
        },
        {
            id: gameState.playerTwo.id,
            name: gameState.playerTwo.name,
            position: p2Position
        }
    ]

    const activeWalls = translateWallsToClient(gameState.walls || []);
    const winnerPlayerName = getPlayerById(gameState, lobbyContext.winnerPlayerId || -1)?.name || "Unknown";

    return (
        <div className="game-frame">
            <WallPicker walls={gameState.playerOne.wallsRemaining} position="top"/>
            <Board players={players}
                    requestPlayerMove={requestPlayerMove}
                    requestWallPlacement={requestWallPlacement}
                    activeWalls={activeWalls}
                    currentTurnPlayerId={gameState.currentTurnPlayerId}
            />
            <WallPicker walls={gameState.playerTwo.wallsRemaining} position="bottom"/>
            <GameOverDialog winnerPlayerName={winnerPlayerName}/>
        </div>
    )
}