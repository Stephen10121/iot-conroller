const WIN_URL = new URL(window.location.href);
const SOCKET_ADDRESS = (WIN_URL.protocol.includes("https") ? "wss://" : "ws://") + WIN_URL.host + "/socket";
const indicator = document.querySelector("#socketIndicator")
const indicatorText = document.querySelector("#socketIndicatorText")
let ws;

function setStatus(status) {
    if (!indicator || !indicatorText) return
    switch (status) {
        case "connected":
            indicator.style.backgroundColor = "#DBF4A7";
            indicatorText.innerText = "Connected";
            break
        case "connecting":
            indicator.style.backgroundColor = "#ed8e1d";
            indicatorText.innerText = "Reconnecting";
            break
        default:
            indicator.style.backgroundColor = "#db4b4b";
            indicatorText.innerText = "Disconnected";
    }
}

function connect() {
    ws = new WebSocket(SOCKET_ADDRESS);

    function socketOpen() {
        // console.log(event);
        console.log("Socket connected.")
        setStatus("connected");
        ws.send("test");
    }

    function socketMessage(event) {
        console.log(JSON.parse(event.data));
        // console.log(event.data);
    }

    ws.onopen = socketOpen;
    ws.onmessage = socketMessage;

    ws.onclose = function(e) {
        if (indicatorText.innerText === "Connected") {
            setStatus("disconnected");
        } else {
            setStatus("connecting");
        }
        console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(function() {
            connect();
        }, 1000);
    };

    ws.onerror = function(err) {
        // indicator.style.backgroundColor = "#d13a3a";
        // indicatorText.innerText = "Disconnected";
        setStatus();
        console.error('Socket encountered error: ', err.message, 'Closing socket');
        ws.close();
    };
}
  
connect();

function makeTestCommand() {
    ws.send("make test:setlight500:1:topic/testPong");
}

function makeAnotherTestCommand() {
    ws.send("make test2:setlight500:1:topic/testPong");
}

function runTestCommand() {
    ws.send("test ttesting something");
}

function deleteTestCommand() {
    ws.send("delete test");
}

function getTestCommand() {
    ws.send("get");
}