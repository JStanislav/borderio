import type { Player } from "../../game/player";
import type { Lobby } from "../../game/lobby/lobby";
import type { MatchConfiguration } from "../../game/MatchConfiguration.ts"

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
