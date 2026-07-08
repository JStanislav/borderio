import { useContext } from "react"
import { LobbyContext } from "../../App.tsx"
import type { LobbyPlayer } from "../../game/lobby/lobby"
import type { MatchConfiguration } from "../../game/MatchConfiguration"
import { PlayerContext } from "../../App.tsx"
import { canDisplayStartButton } from "./lobby.ts"
import "./lobby.css"
import { PlayerCard } from "./PlayerCard.tsx"

interface Props {
    matchConfiguration: MatchConfiguration
    players: LobbyPlayer[]

    actions: Actions
}

interface Actions {
    onPlayerClickStartGame: () => void
    toggleReady: () => void
}

export const Lobby = (props: Props) => {
    const lobby = useContext(LobbyContext);
    const player = useContext(PlayerContext); 

    return ( 
        <div className="lobby-container">
            <div className="lobby-content">
                <div className="lobby-players">
                    {props.players.map((lobbyPlayer) => <PlayerCard name={lobbyPlayer.name} isReady={lobbyPlayer.ready} isHost={lobbyPlayer.host} />)}
                </div>
                <div className="lobby-actions">
                    <button onClick={props.actions.toggleReady}>{player.ready ? "Unready" : "Ready"}</button>
                    {canDisplayStartButton(lobby, props.matchConfiguration, player) && <button onClick={props.actions.onPlayerClickStartGame}>Start</button>}
                </div>
                <div className="lobby-info">
                    <div>Game ID: {lobby.id}</div>
                    <div>Players: {lobby.players.length}/{props.matchConfiguration.playerAmount}</div>
                </div>
            </div>
        </div>
    )
}