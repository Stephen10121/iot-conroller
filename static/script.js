const WIN_URL = new URL(window.location.href);
const SOCKET_ADDRESS = (WIN_URL.protocol.includes("https") ? "wss://" : "ws://") + WIN_URL.host + "/socket";

const websocket = new WebSocket(SOCKET_ADDRESS);

websocket.onopen = (event) => {
    // console.log(event);
    websocket.send("test");
}

websocket.onmessage = (event) => {
    console.log(JSON.parse(event.data));
}