import "./lobby.css"

interface Props {
    name: string,
    isReady: boolean,
    isHost: boolean
}

export const PlayerCard = (props: Props) => {
    return (
        <div key={`${props.name}-card`} className="player-card-container">
            <div className="player-card-leftside">
                <div className="player-card-name">{props.name}</div>
                {props.isHost && <div className="player-card-host">Host</div>}
            </div>
            <div className={props.isReady ? "player-card-ready" : "player-card-not-ready"}>{props.isReady ? "Ready" : "Not Ready"}</div>
        </div>
    )
}