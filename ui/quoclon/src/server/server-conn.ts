import {config} from "../../config/config"
const wsURI = `${config.wsProtocol}${config.serverWSUrl}`;
let websocket: WebSocket | null = null

// The server keeps the player in the game for a grace period after a disconnection,
// so an unexpected drop (e.g. network issues) is retried a few times within that period.
const reconnectDelayMs = 3000;
const maxReconnectAttempts = 15;

export const connect = (
                hash: string,
                ppid: string,
                name: string,
                onMessage: (ev: MessageEvent) => void,
                redirectToHome: () => void,
                reconnectAttempt: number = 0,
            ) => {
    let wasConnected = false;
    const socket = new WebSocket(`${wsURI}/play/${hash}?ppid=${ppid}&name=${name}`);
    websocket = socket;

    socket.addEventListener("open", () => {
        wasConnected = true;
        reconnectAttempt = 0;
    });

    socket.addEventListener("error", (event) => {
        if (!wasConnected && reconnectAttempt === 0 && (socket.readyState === 3)) {
            // the websocket couldn't connect in the first place
            redirectToHome()
        }
        console.error("WebSocket error occurred", event);
    })

    socket.addEventListener("message", (ev) => {
        onMessage(ev)
    })

    socket.addEventListener("close", (ev: CloseEvent) => {
        console.warn("disconnected, reason: ", ev.reason, "code: ", ev.code);

        // Closed on purpose (by us or by the server), or a newer connection replaced this one
        if (ev.code === 1000 || websocket !== socket) {
            return;
        }

        if (!wasConnected && reconnectAttempt === 0) {
            return;
        }

        if (reconnectAttempt >= maxReconnectAttempts) {
            redirectToHome();
            return;
        }

        setTimeout(() => {
            if (websocket !== socket) {
                return;
            }
            connect(hash, ppid, name, onMessage, redirectToHome, reconnectAttempt + 1);
        }, reconnectDelayMs);
    })
}

export const send = (name:string, data:any) => {
    if (websocket != null && websocket.readyState === WebSocket.OPEN) {
        websocket.send(JSON.stringify({type: name, ...data}));
    }
}

export const closeConn = () => {
    if (websocket != null) {
        const socket = websocket;
        // cleared first so the close listener doesn't try to reconnect
        websocket = null;
        socket.close();
    }
}

export const gracefullyCloseConnection = (reason: string) => {
    if (websocket != null) {
        const socket = websocket;
        // cleared first so the close listener doesn't try to reconnect
        websocket = null;
        socket.close(1000, reason);
    }
}
