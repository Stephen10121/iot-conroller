package function

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

func RunCommand(command []string, params Command, conn *websocket.Conn) error {
	joinedCommand := strings.Join(command, " ")
	if params.RelayArguments {
		if len(command) != 2 {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"no arguments given","command":"`+joinedCommand+`"}`))

			if err != nil {
				return err
			}
			return nil
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command":"`+joinedCommand+`","data":"sent command"}`))
		fmt.Println("Running command: " + command[0] + ". With args: " + command[1])
		if err != nil {
			return err
		}
		return nil
	}
	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command":"`+joinedCommand+`","data":"sent command"}`))
	fmt.Println("Running command: " + command[0])

	if err != nil {
		return err
	}
	return nil
}

func CustomCommandChecker(str string, commands map[string]Command) bool {
	_, ok := commands[str]

	return ok
}
