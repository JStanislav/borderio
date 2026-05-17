import type { Player } from "./game/player";
import type { LobbyPlayer } from "./server/messages";


export const generatePPID = () => {
    const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    let result = "";
    const charactersLength = characters.length;
    for (let i = 0; i < 5; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}

export function canDisplayStartButton(lobbyPlayers: LobbyPlayer[], player: Player) {
    return (player.id === 1 && lobbyPlayers.every((lobbyPlayer) => lobbyPlayer.ready))
}