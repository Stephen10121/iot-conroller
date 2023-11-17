const WIN_URL = new URL(window.location.href);
const SOCKET_ADDRESS = (WIN_URL.protocol.includes("https") ? "wss://" : "ws://") + WIN_URL.host + "/socket";

function connect() {
    let ws = new WebSocket(SOCKET_ADDRESS);

    function socketOpen() {
        // console.log(event);
        console.log("Socket connected.")
        ws.send("test");
    }

    function socketMessage(event) {
        console.log(JSON.parse(event.data));
    }

    ws.onopen = socketOpen;
    ws.onmessage = socketMessage;

    ws.onclose = function(e) {
        console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(function() {
            connect();
        }, 1000);
    };

    ws.onerror = function(err) {
        console.error('Socket encountered error: ', err.message, 'Closing socket');
        ws.close();
    };
}
  
connect();