package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID             string
	ConversationID string
	ProfileID      string
	Conn           *websocket.Conn
	Send           chan []byte
	Hub            *Hub
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("Invalid message:", err)
			continue
		}
		msg.SenderID = c.ProfileID

		if msg.ConversationID == "" {
			msg.ConversationID = c.ConversationID
		}

		c.Hub.Broadcast <- msg
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
