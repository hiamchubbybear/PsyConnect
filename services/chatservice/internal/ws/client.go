package ws

import (
	"encoding/json"
	"fmt"

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
			break
		}
		var incoming map[string]interface{}
		if err := json.Unmarshal(raw, &incoming); err != nil {
			if incoming["_system"] == true {
				continue
			}
			fmt.Println("Invalid JSON payload:", string(raw))
			continue
		}

		text, _ := incoming["text"].(string)
		c.Hub.Broadcast <- Message{
			ConversationID: c.ConversationID,
			SenderID:       c.ProfileID,
			Content:        []byte(text),
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
