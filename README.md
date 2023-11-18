# Iot Controller
This is essentially a relay to a broker service that connect to multiple iot devices. This project provides a web gui to help view, make, control, and delete iot devices.
## Usage
To make a custom command, use the make command in when connected to the websocket.
```javascript
// Each configuration parameter is separated by a colon.
ws.send("make testCommand:commandThatsSentToBroker:sendArgsToBrokerBool:callback");

// For example
ws.send("make test:setlight:1:lightIsSet")
```

To use the custom command, send this to the websocket:

```javascript
ws.send("testCommand (if you set sendArgsToBrokerBool to 1 add arguments here.)");

// For example
ws.send("test brightness:500")
```

If you set the ```sendArgsToBrokerBool``` to 1, any arguments you pass in this command, will automatically get sent to the broker itself.