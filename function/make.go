package function

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

func Make(str []string, conn *websocket.Conn, commands map[string]string) error {
	if len(str) > 1 {
		fmt.Println("Making command: " + str[1])

		newCommandSplit := strings.Split(str[1], ":")
		if len(newCommandSplit) != 2 {
			err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"argument given is not formatted properly","command":"`+strings.Join(str, " ")+`"}`))

			if err != nil {
				return err
			}

			return nil
		}

		commands[newCommandSplit[0]] = newCommandSplit[1]
		fmt.Println(commands)

		err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command":"`+strings.Join(str, " ")+`","data":"made command"}`))

		if err != nil {
			return err
		}

		return nil
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"no arguments given","command":"`+str[0]+`"}`))

	if err != nil {
		return err
	}

	return nil
}
