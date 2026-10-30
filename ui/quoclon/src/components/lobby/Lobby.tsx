import { useContext, useEffect, useState } from "react"
import { LobbyContext } from "../../App.tsx"
import type { LobbyPlayer } from "../../game/lobby/lobby"
import type { MatchConfiguration } from "../../game/MatchConfiguration"
import { canDisplayStartButton } from "./lobby.ts"
import "./lobby.css"
import { PlayerCard } from "./PlayerCard.tsx"
import { useAuth } from "../../contexts/auth-provider.tsx"
import { SpectatorsSection } from "./SpectatorsSection.tsx"

interface Props {
    matchConfiguration: MatchConfiguration
    players: LobbyPlayer[]

    actions: Actions
}

interface Actions {
    onPlayerClickStartGame: (ppid: string) => void
    toggleReady: (ppid: string) => void
}

export const Lobby = (props: Props) => {
    const lobby = useContext(LobbyContext);
    const [players, setPlayers] = useState<LobbyPlayer[]>([]);
    const [spectators, setSpectators] = useState<LobbyPlayer[]>([]);
    const { user } = useAuth();

    useEffect( () => {
        setPlayers(props.players.filter(p => !p.spectator));
        setSpectators(props.players.filter(p => p.spectator));
    }, [props.players]);

    if (!user) {
        return null;
    }

    const player = props.players.find(p => p.id === user.id);
    if (player === undefined) {
        return null;
    }

    return ( 
        <div className="lobby-container">
            <div className="lobby-content">
                <div className="lobby-players">
                    <div>
                        {players.map((lobbyPlayer) => <PlayerCard key={`player-${lobbyPlayer.id}-card`} name={lobbyPlayer.name} isReady={lobbyPlayer.ready} isHost={lobbyPlayer.host} />)}
                    </div>
                    <SpectatorsSection spectators={spectators.map((s) => s.name)} />
                    {/* <SpectatorsSection spectators={lobby.players.filter(p => p.spectator).map(p => p.name)}/> */}
                </div>
                <div className="lobby-actions">
                    <button 
                        className="action-button" 
                        style={{visibility: player.spectator ? "hidden" : "visible"}} 
                        onClick={() => props.actions.toggleReady(user.ppid)}>
                            {player.ready ? "Unready" : "Ready"}
                    </button>
                    <button className={`action-button ${!canDisplayStartButton(lobby, props.matchConfiguration, player) ? "no-visible" : ""}`} onClick={() => props.actions.onPlayerClickStartGame(user.ppid)}>Start</button>
                </div>
                <div className="lobby-info">
                    <div>Game ID: {lobby.id}</div>
                    <div>Players: {players.length}/{props.matchConfiguration.playerAmount}</div>
                </div>
            </div>
        </div>
    )
}