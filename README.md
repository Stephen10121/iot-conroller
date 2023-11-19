# Iot Controller
This is essentially a relay to a broker service that connect to multiple iot devices. This project provides a web gui to help view, make, control, and delete iot devices.
## Usage
To make a custom command, use the make command in when connected to the websocket.
```javascript
// Each configuration parameter is separated by a colon.
ws.send("make customCommand:publishTopic:subscribeTopic");

// For example
ws.send("make ledBrightness:topic/led:topic/ledResponse");
```

To use the custom command, send this to the websocket:

```javascript
ws.send("testCommand messageToPublish");

// For example
ws.send("ledBrightness 500")
```