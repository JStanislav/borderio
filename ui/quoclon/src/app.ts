import type { Player } from "./game/player";
import type { Lobby } from "./game/lobby/lobby.ts";
import type { MatchConfiguration } from "./game/MatchConfiguration.ts"


export const generatePPID = () => {
    const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    let result = "";
    const charactersLength = characters.length;
    for (let i = 0; i < 5; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}

export function canDisplayStartButton(lobby: Lobby, matchConfiguration: MatchConfiguration, player: Player) {
    let isHost = false;
    let readyPlayers = 0;
    lobby.players.forEach(p => {
        if (p.host && (p.id === player.id)) {
            isHost = true;
        }
        if (p.ready) {
            readyPlayers++;
        }
    })
    return (isHost && (readyPlayers === matchConfiguration.playerAmount))
}
