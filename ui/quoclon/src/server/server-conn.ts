import {config} from "../../config/config"
const wsURI = `ws://${config.serverUrl}`;
let websocket: WebSocket

export const connect = (onMessage: (ev: MessageEvent) => void) => {
    websocket = new WebSocket(wsURI);
    
    websocket.addEventListener("open", () => {
        console.log("connected");
    }); 

    websocket.addEventListener("error", () => {
        console.log("error");
    })

    websocket.addEventListener("message", (ev) => {
        onMessage(ev)
    })
}

export const send = (name:string, data:any) => {
    if (websocket != null) {
        websocket.send(JSON.stringify({type: name, ...data}));
    }
}

