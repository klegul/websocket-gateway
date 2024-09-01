package ws

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	Clients    map[string]map[*Client]bool
	Broadcast  chan BroadcastMessage
	Register   chan Subscription
	Unregister chan Subscription
}

type Subscription struct {
	Channel string
	Client  *Client
}

type BroadcastMessage struct {
	Channel string
	Message []byte
}
