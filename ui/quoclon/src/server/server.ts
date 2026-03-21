import type { GameState } from "../game/GameState";
import { connect } from "./server-conn";


const onMessage = (ev: MessageEvent, setGameState: (gameState: GameState) => void) => {
    console.log("message arrived");
    console.log("ev", ev.data)
    setGameState(JSON.parse(ev.data));
}


export const startConnection = (setGameState : (gameState: GameState) => void) => {
    // starts socket connection
    connect((ev: MessageEvent) => onMessage(ev, setGameState));
}