import { WallPicker } from "../board/WallPicker"
import "./gameframe.css"
import type { GameState, Player } from "../../game/GameState"
import { requestPlayerMove, requestWallPlacement } from "../../server/server"
import { translateGridPositionToClient, translateWallsToClient } from "../../server/utils"
import { Board } from "../board/Board"
import { useContext, useEffect, useState } from "react"
import { LobbyContext } from "../../App.tsx"


export const GameFrame = ({ gameState }: { gameState: GameState }) => {
    const lobbyContext = useContext(LobbyContext);
    const [p1, setP1] = useState<Player>();
    const [p2, setP2] = useState<Player>();

    useEffect(() => {
        if (!(lobbyContext.players.length === 2)) {
            return;
        }

        setP1(gameState.players.find(p => p.id === lobbyContext.players[0].id));
        setP2(gameState.players.find(p => p.id === lobbyContext.players[1].id));

    }, gameState.players)

    
    if (!p1 || !p2) {
        return
    }

    const p1Position = translateGridPositionToClient(p1.position.row, p1.position.col);
    const p2Position = translateGridPositionToClient(p2.position.row, p2.position.col);
    const players = [
        {
            id: p1.id,
            name: p1.name,
            position: p1Position,
        },
        {
            id: p2.id,
            name: p2.name,
            position: p2Position
        }
    ]

    const activeWalls = translateWallsToClient(gameState.walls || []);

    return (
        <div className="game-frame">
            <WallPicker walls={p1.wallsRemaining} position="top"/> {!lobbyContext.players.find(p => p.id === p1.id)?.connected && "dc"}
            <Board players={players}
                    requestPlayerMove={requestPlayerMove}
                    requestWallPlacement={requestWallPlacement}
                    activeWalls={activeWalls}
                    currentTurnPlayerId={gameState.currentTurnPlayerId}
                    gameOver={lobbyContext.winnerPlayerId !== undefined}
            />
            <WallPicker walls={p2.wallsRemaining} position="bottom"/> {!lobbyContext.players.find(p => p.id === p2.id)?.connected && "dc"}
        </div>
    )
}