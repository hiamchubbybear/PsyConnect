package model

type StartCallPayload struct {
	SessionID      string `json:"sessionId"`
	ConversationID string `json:"conversationId"`
	CallerID       string `json:"callerId"`
	CallerName     string `json:"callerName"`
	ReceiverID     string `json:"receiverId"` 
}
