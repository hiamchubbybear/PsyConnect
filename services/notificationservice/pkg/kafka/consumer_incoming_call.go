package kafka

import (
	"encoding/json"
)

func (c *Consumer) handleIncomingCall(data []byte) error {
	// Define temporary struct to match the event format from Consultation Service
	type IncomingCallEvent struct {
		ReceiverID string      `json:"receiver_id"`
		Payload    interface{} `json:"payload"`
	}

	var event IncomingCallEvent

	// Attempt to unmarshal as the "Wrapper" first, as Consultation Service sends { Key, Value: { Type, Payload } }
	// But wait, the Reader reads msg.Value.
	// In Consultation Producer: p.SendNotification(string(data)).
	// data is json.Marshal(event) where event is the Wrapper.
	// So msg.Value IS the Wrapper.

	type Wrapper struct {
		Key   string `json:"key"`
		Value struct {
			Type    string            `json:"type"`
			Payload IncomingCallEvent `json:"payload"`
		} `json:"value"`
	}

	var w Wrapper
	if err := json.Unmarshal(data, &w); err != nil {
		// Fallback or log error
		// It might be possible that we receive clean payload if using different producer method, but currently it sends wrapper.
		return err
	}

	// Extract the actual payload
	targetPayload := w.Value.Payload

	// Now call the service
	return c.notifSvc.HandleIncomingCall(targetPayload.ReceiverID, targetPayload.Payload)
}
