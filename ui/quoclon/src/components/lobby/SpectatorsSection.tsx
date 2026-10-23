import { useEffect, useState } from "react";
import "./lobby.css"
import { Tooltip } from "react-tooltip";

interface SpectatorsSectionProps {
    spectators: string[]
}

const MaxSpectatorsToShow = 9;


// SpectatorsToDisplay is amount of spectators that can be showed in the spectators section. 
// If there are more than MaxSpectatorsToShow in the lobby, SpectatorsToDisplay will be MaxSpectatorsToShow - 1

export const SpectatorsSection = ({spectators}: SpectatorsSectionProps, ) => {
    const [moreSpectators, setMoreSpectators] = useState<Array<string>>([]);
    const [spectatorsToDisplay, setSpectatorsToDisplay] = useState<Array<string>>(spectators.slice(0, MaxSpectatorsToShow));

    useEffect(() => {
        if (spectators.length > MaxSpectatorsToShow) {
            setMoreSpectators(spectators.slice(MaxSpectatorsToShow));
            setSpectatorsToDisplay(spectators.slice(0, MaxSpectatorsToShow - 1));
        } else{
            setSpectatorsToDisplay(spectators);
            setMoreSpectators([]);
        }
        console.log("SpectatorsSection useEffect called", moreSpectators);

        if (spectators.length > MaxSpectatorsToShow) {

        }
    }, [spectators]);

    return (
        <div className="lobby-spectators">
            <div className="lobby-spectators-title">Spectators</div>
            <div className="lobby-spectators-list">
                {spectatorsToDisplay.map((spectator) =>
                   <div key={`spectator-${spectator}-card`} className="lobby-spectator-card lobby-spectator-name">
                        <span className="spectator-dot"></span>
                        <span>{spectator}</span>
                    </div>
                )}
                {moreSpectators.length > 0 && <div id="more-spectators" style={{userSelect: "none"}}>
                        +<span>{moreSpectators.length}</span>
                    </div>
                    }
                <Tooltip anchorSelect="#more-spectators" place="right">
                    {moreSpectators.map((spectator) =>
                        <div key={`spectator-${spectator}-card`} className="lobby-spectator-card lobby-spectator-name">
                            <span className="spectator-dot"></span>
                            <span>{spectator}</span>
                        </div>
                    )}
                </Tooltip>
            </div>
        </div>
    )
}