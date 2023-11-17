package function

import (
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	messageqeue "github.com/stephen10121/iot-conroller/messageQeue"
)

type Command struct {
	MqttBrokerId       string `json:"mqttBrokerId"`
	ResponseCallbackId string `json:"responseCallbackId"`
	RelayArguments     bool   `json:"relayArguments"`
	MainCommand        string `json:"mainCommand"`
	MqttCommand        string `json:"mqttCommand"`
}

func Make(str string, conn *websocket.Conn, commands map[string]Command, client mqtt.Client) error {
	if len(str) > 1 {
		fmt.Println("Making command: " + str)

		newCommandSplit := strings.Split(str, ":")
		if len(newCommandSplit) != 5 {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"argument given is not formatted properly","command":"`+str+`"}`))

			if err != nil {
				return err
			}

			return nil
		}

		if newCommandSplit[0] == "make" {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"command is reserved","command":"`+str+`"}`))

			if err != nil {
				return err
			}

			return nil
		}

		_, ok := commands[newCommandSplit[0]]
		if ok {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"command already in use","command":"`+str+`"}`))

			if err != nil {
				return err
			}

			return nil
		}

		commands[newCommandSplit[0]] = Command{
			MqttCommand:        newCommandSplit[0],
			MainCommand:        newCommandSplit[1],
			RelayArguments:     newCommandSplit[2] == "1",
			MqttBrokerId:       newCommandSplit[3],
			ResponseCallbackId: newCommandSplit[4],
		}
		a, _ := json.Marshal(commands[newCommandSplit[0]])

		err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command":"`+str+`","data":`+string(a)+`}`))
		messageqeue.Subscribe(client, newCommandSplit[4])

		if err != nil {
			return err
		}

		return nil
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"no arguments given","command":"`+str+`"}`))

	if err != nil {
		return err
	}

	return nil
}
