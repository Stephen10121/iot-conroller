package function

import (
	"fmt"

	"github.com/gorilla/websocket"
	messageqeue "github.com/stephen10121/iot-conroller/messageQeue"
)

func Delete(str string, conn *websocket.Conn, commands map[string]Command) error {
	if len(str) == 0 {
		err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"argument given is not formatted properly","command":"`+str+`"}`))

		if err != nil {
			return err
		}

		return nil
	}
	fmt.Println("Deleting command: " + str)

	messageqeue.Unsubscribe(commands[str].ResponseCallbackId)
	delete(commands, str)

	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command": "delete `+str+`","data":"deleted command"}`))

	if err != nil {
		return err
	}

	return nil
}
