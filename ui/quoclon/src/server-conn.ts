const wsURI = "ws://localhost:8080";
let websocket: WebSocket

export const connect = () => {
    websocket = new WebSocket(wsURI);
    
    websocket.addEventListener("open", () => {
        console.log("connected");
    }); 

    websocket.addEventListener("error", () => {
        console.log("error");
    })

    websocket.addEventListener("message", (ev) => {
        console.log(ev.data);
    })
}

export const send = (name:string, data:string) => {
    if (websocket != null) {
        websocket.send(JSON.stringify({name, data}));
    }
}

