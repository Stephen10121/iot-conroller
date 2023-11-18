package config

// To allow all incoming origins, remove all the strings in the AllowOrigins slice
var AllowOrigins = []string{"localhost:5000"}

// The ip address of the mqtt broker.
var BrokerAddress = "192.168.0.27"

// The port of the mqtt broker.
var BrokerPort = 1883

// How the mqtt broker will id the client
var BrokerClientID = "go_mqtt_client"

// Leave empty if the mqtt broker doesnt require credentials
var BrokerUsername = ""

// Leave empty if the mqtt broker doesnt require credentials
var BrokerPassword = ""

// This is the port of the websocket server and the main web page
var ServerPort = ":5000"
