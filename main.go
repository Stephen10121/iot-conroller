package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/stephen10121/iot-conroller/function"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var commands = make(map[string]function.Command)
var connections = make(map[string]*websocket.Conn)

func removeConnector(conn *websocket.Conn) {
	delete(connections, conn.RemoteAddr().String())
	conn.Close()
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer removeConnector(conn)

	if err != nil {
		log.Println(err)
		return
	}

	connectionId := conn.RemoteAddr().String()
	connections[connectionId] = conn

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			break
		}

		str := strings.Split(string(message), " ")

		switch str[0] {
		case "make":
			err := function.Make(strings.Join(str[1:], " "), conn, commands)

			if err != nil {
				log.Println(err)
			}
			break
		default:
			commandExists := function.CustomCommandChecker(str[0], commands)

			if !commandExists {
				err = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"no such command","command":"`+str[0]+`"}`))

				if err != nil {
					log.Println(err)
				}
				break
			} else {
				function.RunCommand(string(message), commands[str[0]])
			}
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))

	http.HandleFunc("/socket", websocketHandler)
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
