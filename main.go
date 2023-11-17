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

var commands = make(map[string]string)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			break
		}

		str := strings.Split(string(message), " ")

		switch str[0] {
		case "make":
			err := function.Make(str, conn, commands)

			if err != nil {
				log.Println(err)
			}
			break
		default:
			err = conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"no such command","command":"`+str[0]+`"}`))

			if err != nil {
				log.Println(err)
			}
			break
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))

	http.HandleFunc("/socket", websocketHandler)
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
