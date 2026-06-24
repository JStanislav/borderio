import {config} from "../../config/config"
const wsURI = `ws://${config.serverUrl}`;
let websocket: WebSocket

export type actionType = "create" | "join"

export const connect = (
                hash: string,
                actionType: actionType,
                ppid: string,
                onMessage: (ev: MessageEvent) => void,
                redirectToHome: () => void,
            ) => {
    let wasConnected = false;
    websocket = new WebSocket(`${wsURI}/${hash}?action=${actionType}&ppid=${ppid}`);
    
    websocket.addEventListener("open", () => {
        console.log("connected");
        wasConnected = true;
    }); 

    websocket.addEventListener("error", (event) => {
        if (!wasConnected && (websocket.readyState === 3)) {
            // the websocket couldn't connect in the first place
            redirectToHome()
        }
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

export const closeConn = () => {
    if (websocket != null) {
        websocket.close();
    }
}

