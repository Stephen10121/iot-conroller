package function

import (
	"encoding/json"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	messageqeue "github.com/stephen10121/iot-conroller/messageQeue"
)

type Command struct {
	ResponseCallbackId string `json:"responseCallbackId"`
	PublishTopic       string `json:"publishTopic"`
}

func Make(str string, conn *websocket.Conn, commands map[string]Command, client mqtt.Client, connections map[string]*websocket.Conn) error {
	if len(str) > 0 {
		log.Println("Making command: " + str)

		newCommandSplit := strings.Split(str, ":")
		if len(newCommandSplit) != 3 {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"argument given is not formatted properly","command":"`+str+`","type":"makeCommandCallback"}`))

			if err != nil {
				return err
			}

			return nil
		}

		if newCommandSplit[0] == "make" || newCommandSplit[0] == "delete" || newCommandSplit[0] == "get" {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"command is reserved","command":"`+str+`","type":"makeCommandCallback"}`))

			if err != nil {
				return err
			}

			return nil
		}

		_, ok := commands[newCommandSplit[0]]
		if ok {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"command already in use","command":"`+str+`","type":"makeCommandCallback"}`))

			if err != nil {
				return err
			}

			return nil
		}

		commands[newCommandSplit[0]] = Command{
			PublishTopic:       newCommandSplit[1],
			ResponseCallbackId: newCommandSplit[2],
		}
		a, _ := json.Marshal(commands[newCommandSplit[0]])

		err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command":"`+str+`","data":`+string(a)+`,"type":"makeCommandCallback"}`))
		messageqeue.Subscribe(client, newCommandSplit[2], connections)

		if err != nil {
			return err
		}

		return nil
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"no arguments given","command":"`+str+`","type":"makeCommandCallback"}`))

	if err != nil {
		return err
	}

	return nil
}
