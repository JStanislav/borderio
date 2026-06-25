import toast from "react-hot-toast";
import { type GameState } from "../game/GameState";
import { connect,send, type actionType } from "./server-conn";
import { translateGridPositionToServer, translateWallGridPositionToServer } from "./utils";
import type { IMessage, LobbyMessage, PlayerConfigurationMessage, PlayerJoinedMessage, WillTimeoutMessage } from "./messages";
import type { Player } from "../game/player";
import { config } from "../../config/config";
import type { Lobby } from "../game/lobby/lobby";
import type { MatchConfiguration } from "../game/MatchConfiguration";

const serverURL = `http://${config.serverUrl}`

const onMessage = (ev: MessageEvent,
                    setGameState: (gameState: GameState) => void,
                    setPlayerConfig: (player: Player) => void,
                    setLobby: (lobby: Lobby) => void,
                    setMatchConfiguration: (matchConfiguration: MatchConfiguration) => void) => {

    console.log("message arrived");
    console.log("ev", ev.data)

    const data = JSON.parse(ev.data) as IMessage<any>;
    
    if (data.type === "gameState") {
        setGameState(data.payload);
    } else if (data.type === "error") {
        toast.error(`Error: ${data.payload}`);
    }
    if (data.type === "playerConfiguration") {
        const config = data.payload as PlayerConfigurationMessage;
        toast.success(`You are ${config.name} (id: ${config.id}) with ppid: ${config.ppid}`);
        setPlayerConfig({ id: config.id, name: config.name, ppid: config.ppid, ready: false });
    }
    if (data.type === "matchConfiguration") {
        const config = data.payload as MatchConfiguration;
        setMatchConfiguration(config);
    }
    if (data.type === "lobby") {
        const config = data.payload as LobbyMessage;
        setLobby({ players: config.players, winnerPlayerId: config.winnerPlayerId, playerAmount: 2 });
    }
    if (data.type === "joined") { 
        const joined = data.payload as PlayerJoinedMessage;
        toast.success(`${joined.name}[${joined.id}] has joined game!`);
    }
    if (data.type === "playerLeft") {
        const left = data.payload as PlayerJoinedMessage;
        toast.success(`${left.name} has left game!`);
    }
    if (data.type === "willTimeOut") {
        const willTimeOut = data.payload as WillTimeoutMessage
        toast(`The lobby will be closed in ${willTimeOut.span}`, {icon: "⌛"});
    }
    if (data.type === "gameTimedOut") {
        toast("The connection has been closed.", {icon: "👋"});
    }
}


export const startConnection = (hash: string,
                                action: actionType,
                                ppid: string,
                                setGameState: (gameState: GameState) => void,
                                setPlayerConfig: (player: Player) => void,
                                setLobby: (lobby: Lobby) => void,
                                setMatchConfiguration: (matchConfiguration: MatchConfiguration) => void,
                                redirectToHome: () => void) => {

    // starts socket connection
   connect(hash, action, ppid, (ev: MessageEvent) => onMessage(ev, setGameState, setPlayerConfig, setLobby, setMatchConfiguration), redirectToHome);
}

export const requestPlayerMove = (ppid: string, row: number, col: number) => {
    const { s_row, s_col } = translateGridPositionToServer(row, col);
    const type = "playerMove";
    const target = { row: s_row,col: s_col };
    send(type, { payload: { target }, ppid});
}

export const requestWallPlacement = (ppid: string, row: number, col: number, orientation: "horizontal" | "vertical") => {
    let wallPositions = translateWallGridPositionToServer(row, col, orientation);
    const type = "wallPlacement";
    const wallTarget = { cellA: { row: wallPositions.cellA.row, col: wallPositions.cellA.column }, cellB: { row: wallPositions.cellB.row, col: wallPositions.cellB.column }, orientation };
    send(type, { payload: { wallTarget: wallTarget }, ppid});
}

export async function GameExist(hash: string) {
    if (hash === "") return false;

    const res = await fetch(`${serverURL}/ping/${hash}`, { method: "GET" });

    if (res.status !== 200) return false
    return true

}