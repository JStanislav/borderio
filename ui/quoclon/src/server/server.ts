import toast from "react-hot-toast";
import { type GameState } from "../game/GameState";
import { connect,send, type actionType } from "./server-conn";
import { translateGridPositionToServer, translateWallGridPositionToServer } from "./utils";
import type { LobbyMessage, LobbyPlayer, PlayerConfigurationMessage, PlayerJoinedMessage } from "./messages";
import type { Player } from "../game/player";
import { config } from "../../config/config";

const serverURL = `http://${config.serverUrl}`

const onMessage = (ev: MessageEvent,
                    setGameState: (gameState: GameState) => void,
                    setPlayerConfig: (player: Player) => void,
                    setLobbyPlayers: (players: LobbyPlayer[]) => void) => {

    console.log("message arrived");
    console.log("ev", ev.data)

    const data = JSON.parse(ev.data);
    
    if (data.type === "gameState") {
        setGameState(data);
    } else if (data.type === "error") {
        toast.error(`Error: ${data.message}`);
    }
    if (data.type === "playerConfiguration") {
        const config = data as PlayerConfigurationMessage;
        toast.success(`You are ${config.name} (id: ${config.id}) with ppid: ${config.ppid}`);
        setPlayerConfig({ id: config.id, name: config.name, ppid: config.ppid, ready: false });
    }
    if (data.type === "lobby") {
        const config = data as LobbyMessage;
        setLobbyPlayers(config.players);
    }
    if (data.type === "joined") { 
        const joined = data as PlayerJoinedMessage;
        toast.success(`${joined.name}[${joined.id}] has joined game!`);
    }
}


export const startConnection = (hash: string,
                                action: actionType,
                                ppid: string,
                                setGameState: (gameState: GameState) => void,
                                setPlayerConfig: (player: Player) => void,
                                setLobbyPlayers: (players: LobbyPlayer[]) => void) => {
    // starts socket connection
    connect(hash, action, ppid, (ev: MessageEvent) => onMessage(ev, setGameState, setPlayerConfig, setLobbyPlayers));
}

export const requestPlayerMove = (playerId: number, row: number, col: number) => {
    const { s_row, s_col } = translateGridPositionToServer(row, col);
    const type = "playerMove";
    const target = { row: s_row,col: s_col };
    send(type, { playerId, target });
}

export const requestWallPlacement = (playerId: number, row: number, col: number, orientation: "horizontal" | "vertical") => {
    let wallPositions = translateWallGridPositionToServer(row, col, orientation);
    const type = "wallPlacement";
    const wallTarget = { cellA: { row: wallPositions.cellA.row, col: wallPositions.cellA.column }, cellB: { row: wallPositions.cellB.row, col: wallPositions.cellB.column }, orientation };
    send(type, { playerId, wallTarget: wallTarget });
}

export async function GameExist(hash: string) {
    if (hash === "") return false;

    const res = await fetch(`${serverURL}/ping/${hash}`, { method: "GET" });

    if (res.status !== 200) return false
    return true

}