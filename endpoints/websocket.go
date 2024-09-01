package endpoints

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"websocket-gateway/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var hub = ws.Hub{
	Clients:    make(map[string]map[*ws.Client]bool),
	Broadcast:  make(chan ws.BroadcastMessage),
	Register:   make(chan ws.Subscription),
	Unregister: make(chan ws.Subscription),
}

func Run() {
	go hub.Run()
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query().Get("channel")
	if channel == "" {
		http.Error(w, "channel required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error during WebSocket upgrade: %v", err)
	}

	client := &ws.Client{Conn: conn, Send: make(chan []byte, 256)}

	go client.WritePump()
	go client.ReadPump(&hub)

	hub.Register <- ws.Subscription{Channel: channel, Client: client}
}
