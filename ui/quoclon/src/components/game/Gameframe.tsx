import { WallPicker } from "../board/WallPicker"
import "./gameframe.css"
import { GetDefaultPlayer, type GameState, type Player } from "../../game/GameState"
import { requestPlayerMove, requestWallPlacement } from "../../server/server"
import { translateGridPositionToClient, translateWallsToClient } from "../../server/utils"
import { Board } from "../board/Board"
import { useContext, useEffect, useState } from "react"
import { LobbyContext } from "../../App.tsx"
import { useAuth } from "../../contexts/auth-provider.tsx"
import { GameTitle } from "./GameTitle.tsx"

export const GameFrame = ({ gameState }: { gameState: GameState }) => {
    const lobbyContext = useContext(LobbyContext);
    const {user} = useAuth()
    const [players, setPlayers] = useState<Player[]>([])

    useEffect(() => {
        const playersNotSpectators = lobbyContext.players.filter(p => !p.spectator)
        if ((playersNotSpectators.length !== lobbyContext.playerAmount)) {
            return
        }

        setPlayers(playersNotSpectators.map(player => {
            let p = gameState.players.get(player.id)
            if (!p) {
                // This should never happen, but if it does, we will create a default player to avoid crashing the game
                // It happens because a race condition where the lobby is updated before the game state is updated,
                // so the player is not yet in the game state when the lobby is updated.
                console.warn("no player")
                p = GetDefaultPlayer(player.id, player.name)
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
            <GameTitle players={players.filter(p => p.id !== -1).map(p => p.name)} currentUser={user?.name || ""} />
            <div className="game-frame-board">
                <WallPicker walls={players[0]?.wallsRemaining || 0} position="top"/>
                {!lobbyContext.players.find(p => p.id === players[0]?.id)?.connected && "dc"}
                <Board players={players}
                        requestPlayerMove={requestPlayerMove}
                        requestWallPlacement={requestWallPlacement}
                        activeWalls={activeWalls}
                        currentTurnPlayerId={gameState.currentTurnPlayerId}
                        gameOver={lobbyContext.winnerPlayerId !== undefined}
                        spectator={user ? user.id === -1 : true}
                />
                <WallPicker walls={players[1]?.wallsRemaining || 0} position="bottom"/>
                {!lobbyContext.players.find(p => p.id === players[1]?.id)?.connected && "dc"}
            </div>
        </div>
    )
}