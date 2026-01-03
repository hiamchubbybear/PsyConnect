package signaling

import (
	"chatservice/internal/ws"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

// WebRTC signaling data structures
type SDPData struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"` // "offer" or "answer"
}

type ICECandidateData struct {
	Candidate     string `json:"candidate"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex int    `json:"sdpMLineIndex"`
}

type CallSession struct {
	SessionID      string     `json:"sessionId"`
	ConversationID string     `json:"conversationId"`
	CallerID       string     `json:"callerId"`
	CalleeID       string     `json:"calleeId"`
	StartTime      time.Time  `json:"startTime"`
	EndTime        *time.Time `json:"endTime,omitempty"`
	Status         string     `json:"status"` // "ringing", "active", "ended"
}

type SignalingService struct {
	hub      *ws.Hub
	sessions map[string]*CallSession // sessionId -> CallSession
}

func NewSignalingService(hub *ws.Hub) *SignalingService {
	service := &SignalingService{
		hub:      hub,
		sessions: make(map[string]*CallSession),
	}

	// Register WebRTC message handlers
	hub.RegisterHandler(ws.MessageTypeOffer, service.HandleOffer)
	hub.RegisterHandler(ws.MessageTypeAnswer, service.HandleAnswer)
	hub.RegisterHandler(ws.MessageTypeICE, service.HandleICECandidate)
	hub.RegisterHandler(ws.MessageTypeLeave, service.HandleLeave)

	log.Println("✅ Signaling service initialized with WebRTC handlers")
	return service
}

// HandleOffer - Caller sends offer to callee
func (s *SignalingService) HandleOffer(hub *ws.Hub, message ws.Message) {
	log.Printf("📞 Handling offer from %s to %s", message.SenderID, message.ReceiverID)

	// Parse SDP offer
	var sdpData SDPData
	if err := json.Unmarshal(message.Data, &sdpData); err != nil {
		log.Printf("❌ Error parsing SDP offer: %v", err)
		return
	}

	// Create call session
	sessionID := uuid.New().String()
	session := &CallSession{
		SessionID:      sessionID,
		ConversationID: message.ConversationID,
		CallerID:       message.SenderID,
		CalleeID:       message.ReceiverID,
		StartTime:      time.Now(),
		Status:         "ringing",
	}
	s.sessions[sessionID] = session

	// Forward offer to callee
	offerMessage := ws.Message{
		Type:           ws.MessageTypeOffer,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Data:           message.Data,
	}

	// Add session ID to data
	dataWithSession := map[string]interface{}{
		"sdp":       sdpData.SDP,
		"type":      sdpData.Type,
		"sessionId": sessionID,
	}
	dataBytes, _ := json.Marshal(dataWithSession)
	offerMessage.Data = dataBytes

	hub.SendToClient(message.ConversationID, message.ReceiverID, offerMessage)
	log.Printf("✅ Offer forwarded to %s, session: %s", message.ReceiverID, sessionID)
}

// HandleAnswer - Callee sends answer to caller
func (s *SignalingService) HandleAnswer(hub *ws.Hub, message ws.Message) {
	log.Printf("📞 Handling answer from %s to %s", message.SenderID, message.ReceiverID)

	// Parse answer data
	var answerData map[string]interface{}
	if err := json.Unmarshal(message.Data, &answerData); err != nil {
		log.Printf("❌ Error parsing answer: %v", err)
		return
	}

	sessionID, ok := answerData["sessionId"].(string)
	if !ok {
		log.Printf("❌ No session ID in answer")
		return
	}

	// Update session status
	if session, exists := s.sessions[sessionID]; exists {
		session.Status = "active"
		log.Printf("✅ Call session %s is now active", sessionID)
	}

	// Forward answer to caller
	answerMessage := ws.Message{
		Type:           ws.MessageTypeAnswer,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Data:           message.Data,
	}

	hub.SendToClient(message.ConversationID, message.ReceiverID, answerMessage)
	log.Printf("✅ Answer forwarded to %s", message.ReceiverID)
}

// HandleICECandidate - Exchange ICE candidates
func (s *SignalingService) HandleICECandidate(hub *ws.Hub, message ws.Message) {
	log.Printf("🧊 Handling ICE candidate from %s to %s", message.SenderID, message.ReceiverID)

	// Forward ICE candidate to peer
	iceMessage := ws.Message{
		Type:           ws.MessageTypeICE,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Data:           message.Data,
	}

	hub.SendToClient(message.ConversationID, message.ReceiverID, iceMessage)
	log.Printf("✅ ICE candidate forwarded to %s", message.ReceiverID)
}

// HandleLeave - End call
func (s *SignalingService) HandleLeave(hub *ws.Hub, message ws.Message) {
	log.Printf("👋 Handling leave from %s", message.SenderID)

	// Parse leave data
	var leaveData map[string]interface{}
	if err := json.Unmarshal(message.Data, &leaveData); err != nil {
		log.Printf("❌ Error parsing leave data: %v", err)
		return
	}

	sessionID, ok := leaveData["sessionId"].(string)
	if !ok {
		log.Printf("⚠️  No session ID in leave message")
	} else {
		// End session
		if session, exists := s.sessions[sessionID]; exists {
			now := time.Now()
			session.EndTime = &now
			session.Status = "ended"
			log.Printf("✅ Call session %s ended", sessionID)
		}
	}

	// Notify peer
	leaveMessage := ws.Message{
		Type:           ws.MessageTypeLeave,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Data:           message.Data,
	}

	hub.SendToClient(message.ConversationID, message.ReceiverID, leaveMessage)
	log.Printf("✅ Leave notification sent to %s", message.ReceiverID)
}

// GetActiveSession - Get active call session for conversation
func (s *SignalingService) GetActiveSession(conversationID string) *CallSession {
	for _, session := range s.sessions {
		if session.ConversationID == conversationID && session.Status == "active" {
			return session
		}
	}
	return nil
}

// GetSession - Get session by ID
func (s *SignalingService) GetSession(sessionID string) *CallSession {
	return s.sessions[sessionID]
}

// CleanupOldSessions - Remove ended sessions older than 1 hour
func (s *SignalingService) CleanupOldSessions() {
	cutoff := time.Now().Add(-1 * time.Hour)
	for id, session := range s.sessions {
		if session.Status == "ended" && session.EndTime != nil && session.EndTime.Before(cutoff) {
			delete(s.sessions, id)
			log.Printf("🧹 Cleaned up old session: %s", id)
		}
	}
}
