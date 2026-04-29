import toast from "react-hot-toast";
import { type GameState } from "../game/GameState";
import { connect,send } from "./server-conn";
import { translateGridPositionToServer, translateWallGridPositionToServer } from "./utils";


const onMessage = (ev: MessageEvent, setGameState: (gameState: GameState) => void) => {
    console.log("message arrived");
    console.log("ev", ev.data)
    const data = JSON.parse(ev.data);
    if (data.type === "gameState") {
        setGameState(data);
    } else if (data.type === "error") {
        toast.error(`Error: ${data.message}`);
    }
}


export const startConnection = (setGameState : (gameState: GameState) => void) => {
    // starts socket connection
    connect("testHash", (ev: MessageEvent) => onMessage(ev, setGameState));
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