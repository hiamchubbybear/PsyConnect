package ws

import (
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
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		c.Hub.Broadcast <- Message{
			ConversationID: c.ConversationID,
			SenderID:       c.ID,
			Content:        message,
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
