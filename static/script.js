const WIN_URL = new URL(window.location.href);
const SOCKET_ADDRESS = (WIN_URL.protocol.includes("https") ? "wss://" : "ws://") + WIN_URL.host + "/socket";
const indicator = document.querySelector("#socketIndicator")
let ws;

function connect() {
    ws = new WebSocket(SOCKET_ADDRESS);

    function socketOpen() {
        // console.log(event);
        console.log("Socket connected.")
        indicator.style.backgroundColor = "#07ff07";
        indicator.children[0].innerText = "Connected";
        ws.send("test");
    }

    function socketMessage(event) {
        console.log(JSON.parse(event.data));
    }

    ws.onopen = socketOpen;
    ws.onmessage = socketMessage;

    ws.onclose = function(e) {
        if (indicator.children[0].innerText === "Connected") {
            indicator.style.backgroundColor = "#d13a3a";
            indicator.children[0].innerText = "Disconnected";
        } else {
            indicator.style.backgroundColor = "#ed8e1d";
            indicator.children[0].innerText = "Reconnecting";
        }
        console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(function() {
            connect();
        }, 1000);
    };

    ws.onerror = function(err) {
        indicator.style.backgroundColor = "#d13a3a";
        indicator.children[0].innerText = "Disconnected";
        console.error('Socket encountered error: ', err.message, 'Closing socket');
        ws.close();
    };
}
  
connect();