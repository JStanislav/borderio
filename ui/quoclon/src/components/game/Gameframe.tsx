import { WallPicker } from "../board/WallPicker"
import "./gameframe.css"
import { type GameState, type Player } from "../../game/GameState"
import { requestPlayerMove, requestWallPlacement } from "../../server/server"
import { translateGridPositionToClient, translateWallsToClient } from "../../server/utils"
import { Board } from "../board/Board"
import { useContext, useEffect, useState } from "react"
import { LobbyContext } from "../../App.tsx"

export const GameFrame = ({ gameState }: { gameState: GameState }) => {
    const lobbyContext = useContext(LobbyContext);
    const [players, setPlayers] = useState<Player[]>([])

    useEffect(() => {
        if ((lobbyContext.players.length !== lobbyContext.playerAmount)) {
            return
        }

        setPlayers(lobbyContext.players.map(player => {
            const p = gameState.players.get(player.id)
            if (!p) {
                throw new Error("unexpected error: player not found")
            }

            return {
                ...p,
                position: translateGridPositionToClient(p.position.row, p.position.col)
            }
        }))

    }, [gameState.players])

    const activeWalls = translateWallsToClient(gameState.walls || []);

    return (
        <div className="game-frame">
            <WallPicker walls={players[0]?.wallsRemaining || 0} position="top"/>
            <Board players={players}
                    requestPlayerMove={requestPlayerMove}
                    requestWallPlacement={requestWallPlacement}
                    activeWalls={activeWalls}
                    currentTurnPlayerId={gameState.currentTurnPlayerId}
                    gameOver={lobbyContext.winnerPlayerId !== undefined}
            />
            <WallPicker walls={players[1]?.wallsRemaining || 0} position="bottom"/>
        </div>
    )
}