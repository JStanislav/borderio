import {config} from "../../config/config"
const wsURI = `ws://${config.serverUrl}`;
let websocket: WebSocket

export const connect = (hash: string, onMessage: (ev: MessageEvent) => void) => {
    websocket = new WebSocket(`${wsURI}/${hash}`);
    
    websocket.addEventListener("open", () => {
        console.log("connected");
    }); 

    websocket.addEventListener("error", (event) => {
        console.error("WebSocket error occurred", event);
    })

    websocket.addEventListener("message", (ev) => {
        onMessage(ev)
    })

    websocket.addEventListener("close", (ev: CloseEvent) => {
        console.log("disconnected, reason: ", ev.reason, "code: ", ev.code);
    })
}

export const send = (name:string, data:any) => {
    if (websocket != null) {
        websocket.send(JSON.stringify({type: name, ...data}));
    }
}

