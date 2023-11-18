package function

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func Get(conn *websocket.Conn, commands map[string]Command) error {
	log.Println("Sending all commands")
	jsonData, err := json.Marshal(commands)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command": "get","data":`+string(jsonData)+`, "type":"getCommandCallback"}`))

	if err != nil {
		return err
	}

	return nil
}
