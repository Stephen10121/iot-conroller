package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/stephen10121/iot-conroller/function"
	messageqeue "github.com/stephen10121/iot-conroller/messageQeue"
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
var messageQeaueConnection mqtt.Client

func removeConnector(conn *websocket.Conn) {
	delete(connections, conn.RemoteAddr().String())
	conn.Close()

	fmt.Println("Closed Connection: " + conn.RemoteAddr().String())
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer removeConnector(conn)

	if err != nil {
		log.Println(err)
		return
	}

	connectionId := conn.RemoteAddr().String()
	fmt.Println("New Connection: " + connectionId)
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
				err := function.RunCommand(str, commands[str[0]], conn, messageQeaueConnection)

				if err != nil {
					log.Println(err)
				}
				break
			}
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	messageQeaueConnection = messageqeue.Initialize()

	http.HandleFunc("/socket", websocketHandler)
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
