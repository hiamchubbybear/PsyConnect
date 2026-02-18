package kafka

import (
	"encoding/json"
	"log"
)

func (c *Consumer) handleIncomingCall(data []byte) error {
	// Define the event structure that matches what Chat Service sends
	type IncomingCallEvent struct {
		EventID   string `json:"eventId"`
		EventType string `json:"eventType"`
		Data      struct {
			SessionID      string `json:"sessionId"`
			ConversationID string `json:"conversationId"`
			CallerID       string `json:"callerId"`
			CallerName     string `json:"callerName"`
			RecipientID    string `json:"recipientId"`
		} `json:"data"`
	}

	var event IncomingCallEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("Failed to unmarshal incoming call event: %v", err)
		return err
	}

	// Log for debugging
	log.Printf("Received incoming call event: %s from %s to %s", event.EventID, event.Data.CallerID, event.Data.RecipientID)

	// Call the notification service to send FCM
	// We pass the whole Data struct as payload, or map it to what HandleIncomingCall expects
	// Assuming HandleIncomingCall takes (receiverID string, payload interface{})
	return c.notifSvc.HandleIncomingCall(event.Data.RecipientID, event.Data)
}
