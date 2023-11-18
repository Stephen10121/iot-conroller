package types

import "github.com/gorilla/websocket"

type ResponderMap struct {
	Command   string
	Recievers []*websocket.Conn
}
