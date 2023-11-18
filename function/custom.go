package function

import (
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	messageqeue "github.com/stephen10121/iot-conroller/messageQeue"
)

func RunCommand(command []string, params Command, conn *websocket.Conn, mqttClient mqtt.Client) error {
	joinedCommand := strings.Join(command, " ")
	actualCommand := params.MainCommand + " " + strings.Join(command[1:], " ")
	if params.RelayArguments {
		if len(command) < 2 {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"no arguments given","command":"`+joinedCommand+`"}`))

			if err != nil {
				return err
			}
			return nil
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command":"`+joinedCommand+`","data":"sent command"}`))
		fmt.Println("Running command: " + command[0] + ". With args: " + command[1] + ". Command sent to broker: " + actualCommand)
		go messageqeue.Publish(mqttClient, actualCommand)

		if err != nil {
			return err
		}
		return nil
	}
	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command":"`+joinedCommand+`","data":"sent command"}`))
	fmt.Println("Running command: " + command[0] + ". Command sent to broker: " + actualCommand)
	go messageqeue.Publish(mqttClient, actualCommand)

	if err != nil {
		return err
	}
	return nil
}

func CustomCommandChecker(str string, commands map[string]Command) bool {
	_, ok := commands[str]

	return ok
}
