package ws

import (
	"encoding/json"
	"log"
	"time"
)

type MessageType string

const (
	MessageTypeChat   MessageType = "chat"
	MessageTypeOffer  MessageType = "offer"
	MessageTypeAnswer MessageType = "answer"
	MessageTypeICE    MessageType = "ice"
	MessageTypeJoin   MessageType = "join"
	MessageTypeLeave  MessageType = "leave"
)

//	type Message struct {
//		Type           MessageType            `json:"type"`
//		ConversationID string                 `json:"conversationId"`
//		SenderID       string                 `json:"senderId"`
//		ReceiverID     string                 `json:"receiverId,omitempty"`
//		Data           map[string]interface{} `json:"data"`
//	}
type Message struct {
	Type           MessageType     `json:"type"`
	ConversationID string          `json:"conversationId"`
	SenderID       string          `json:"senderId"`
	ReceiverID     string          `json:"receiverId"`
	Data           json.RawMessage `json:"data"`
}
type Offer struct {
	Type           MessageType `json:"type"`
	ConversationID string      `json:"conversationId"`
	SenderID       string      `json:"senderId"`
	ReceiverID     string      `json:"receiverId"`
	StartTime      time.Time   `json:"startTime" bson:"startTime"`
	EndTime        time.Time   `json:"endTime" bson:"endTime"`
	SessionId      string      `json:"sessionId" bson:"sessionId"`
}
type MessageHandler func(hub *Hub, message Message)

type Hub struct {
	Clients    map[string][]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message

	handlers map[MessageType]MessageHandler
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
		handlers:   make(map[MessageType]MessageHandler),
	}
}

func (h *Hub) RegisterHandler(msgType MessageType, handler MessageHandler) {
	h.handlers[msgType] = handler
	log.Printf("Registered handler for message type: %s", msgType)
}

func (h *Hub) Run() {
	log.Println("WebSocket Hub started")
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ConversationID] = append(h.Clients[client.ConversationID], client)
			log.Printf("Client registered: %s in conversation: %s", client.ID, client.ConversationID)

		case client := <-h.Unregister:
			clients := h.Clients[client.ConversationID]
			for i, c := range clients {
				if c == client {
					h.Clients[client.ConversationID] = append(clients[:i], clients[i+1:]...)
					close(client.Send)
					log.Printf("Client unregistered: %s from conversation: %s", client.ID, client.ConversationID)
					break
				}
			}

		case message := <-h.Broadcast:
			log.Printf("Broadcasting message type: %s from sender: %s", message.Type, message.SenderID)

			if handler, ok := h.handlers[message.Type]; ok {
				if MessageTypeOffer == message.Type {
					handler(h, message)
				}
				handler(h, message)

			} else {
				log.Printf("No handler registered for message type: %s please change it", message.Type)
			}
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID string, message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	count := 0
	for _, client := range h.Clients[roomID] {
		select {
		case client.Send <- data:
			count++
		default:
			log.Printf("Client send buffer full, skipping: %s", client.ID)
		}
	}
	log.Printf("Broadcasted to %d clients in room: %s", count, roomID)
}

func (h *Hub) SendToClient(roomID, clientID string, message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for _, client := range h.Clients[roomID] {
		if client.ID == clientID {
			select {
			case client.Send <- data:
				log.Printf("Sent message to client: %s", clientID)
			default:
				log.Printf("Client send buffer full: %s", clientID)
			}
			return
		}
	}
	log.Printf("Client not found in room: %s", clientID)
}

func (h *Hub) GetClientsInRoom(roomID string) []*Client {
	return h.Clients[roomID]
}

func (h *Hub) GetClientCount(roomID string) int {
	return len(h.Clients[roomID])
}
