package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

		log.Printf("Recieved message: %s", message)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))

	http.HandleFunc("/socket", websocketHandler)
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
