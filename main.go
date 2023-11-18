package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/stephen10121/iot-conroller/config"
	"github.com/stephen10121/iot-conroller/function"
	messageqeue "github.com/stephen10121/iot-conroller/messageQeue"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		fmt.Println(r.Host)
		if len(config.AllowOrigins) > 0 {
			if slices.Contains(config.AllowOrigins, r.Host) {
				return true
			} else {
				return false
			}
		} else {
			return true
		}
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

	if err != nil {
		log.Println("bob", err)
		return
	}

	defer removeConnector(conn)

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
			err := function.Make(strings.Join(str[1:], " "), conn, commands, messageQeaueConnection, connections)

			if err != nil {
				log.Println(err)
			}
			break
		default:
			if function.CustomCommandChecker(str[0], commands) {
				err := function.RunCommand(str, commands[str[0]], conn, messageQeaueConnection)

				if err != nil {
					log.Println(err)
				}
			} else {
				err = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"no such command","command":"`+str[0]+`"}`))

				if err != nil {
					log.Println(err)
				}

			}
			break
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
