package ws

import (
	"github.com/gorilla/websocket"
	"log"
)

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- Subscription{Channel: "", Client: c}
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf("Received: %s", message)
	}
}

func (c *Client) WritePump() {
	for message := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
	}
	c.Conn.Close()
}

func (hub *Hub) Run() {
	for {
		select {
		case subscription := <-hub.Register:
			if hub.Clients[subscription.Channel] == nil {
				hub.Clients[subscription.Channel] = make(map[*Client]bool)
			}
			hub.Clients[subscription.Channel][subscription.Client] = true
		case subscription := <-hub.Unregister:
			for channel, clients := range hub.Clients {
				if clients[subscription.Client] {
					delete(clients, subscription.Client)
					close(subscription.Client.Send)
					if len(clients) == 0 {
						delete(hub.Clients, channel)
					}
					break
				}
			}
		case message := <-hub.Broadcast:
			clients := hub.Clients[message.Channel]
			for client := range clients {
				select {
				case client.Send <- message.Message:
				default:
					close(client.Send)
					delete(clients, client)
					if len(clients) == 0 {
						delete(hub.Clients, message.Channel)
					}
				}
			}
		}
	}
}
