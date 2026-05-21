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

export function canDisplayStartButton(lobbyPlayers: LobbyPlayer[], amountOfPlayers: number, player: Player) {
    return ((player.id === 1) && lobbyPlayers.length === amountOfPlayers && lobbyPlayers.every((lobbyPlayer) => lobbyPlayer.ready))
}