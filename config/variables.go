package config

// To allow all incoming origins, remove all the strings in the AllowOrigins slice
var AllowOrigins = []string{"localhost:5000"}

// The ip address of the mqtt broker.
var BrokerAddress = "192.168.0.27"

// The port of the mqtt broker.
var BrokerPort = 1883
