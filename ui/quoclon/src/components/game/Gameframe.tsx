import { WallPicker } from "../board/WallPicker"
import "./gameframe.css"
import { getPlayerById, type GameState } from "../../game/GameState"
import { requestPlayerMove, requestWallPlacement } from "../../server/server"
import { translateGridPositionToClient, translateWallsToClient } from "../../server/utils"
import { Board } from "../board/Board"
import { useContext, useEffect, useState } from "react"
import { LobbyContext } from "../../App.tsx"


export const GameFrame = ({ gameState }: { gameState: GameState }) => {
    const [showGameOverDialog, setShowGameOverDialog] = useState(false);
    const lobbyContext = useContext(LobbyContext);

    useEffect(() => {
        if (lobbyContext.winnerPlayerId !== undefined) {
            setShowGameOverDialog(true);
        }
    }, [lobbyContext.winnerPlayerId])

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
            <dialog id="game-over-dialog" open={showGameOverDialog} closedby="any">
                Holi, ganó {getPlayerById(gameState, lobbyContext.winnerPlayerId!)?.name}! Felicidades!
                <button commandfor="game-over-dialog" command="close">Cerrar</button>
            </dialog>
        </div>
    )
}