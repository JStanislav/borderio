import { WallPicker } from "../board/WallPicker"
import "./gameframe.css"
import { getPlayerById, type GameState } from "../../game/GameState"
import { requestPlayerMove, requestWallPlacement } from "../../server/server"
import { translateGridPositionToClient, translateWallsToClient } from "../../server/utils"
import { Board } from "../board/Board"
import { useContext, useEffect, useRef } from "react"
import { LobbyContext, PlayerContext } from "../../App.tsx"
import { getGameOverText } from "./player.type.ts"


export const GameFrame = ({ gameState }: { gameState: GameState }) => {
    const dialogRef = useRef<HTMLDialogElement>(null);
    const lobbyContext = useContext(LobbyContext);
    const playerContext = useContext(PlayerContext);

    useEffect(() => {
        if (lobbyContext.winnerPlayerId !== undefined) {
            onOpenModal();
        }
    }, [lobbyContext.winnerPlayerId])

    const onCloseModal = () => {
        dialogRef.current?.close();
    }
    
    const onOpenModal = () => {
        dialogRef.current?.showModal();
    }

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
            <div className="dialog-container">
                <dialog ref={dialogRef} className="game-over-dialog" closedby="any">
                    <p>{getGameOverText(winnerPlayerName, lobbyContext.winnerPlayerId === playerContext.id)}</p>
                    <button onClick={onCloseModal}>Cerrar</button>
                </dialog>        
            </div>
        </div>
    )
}