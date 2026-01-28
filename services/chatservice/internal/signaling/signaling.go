package signaling

import (
	"chatservice/internal/ws"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SDPData struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type ICECandidateData struct {
	Candidate     string `json:"candidate"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex int    `json:"sdpMLineIndex"`
}

type CallStatus string

const (
	CallStatusRinging CallStatus = "ringing"
	CallStatusActive  CallStatus = "active"
	CallStatusEnded   CallStatus = "ended"
)

type CallSession struct {
	SessionID      string     `json:"sessionId"`
	ConversationID string     `json:"conversationId"`
	CallerID       string     `json:"callerId"`
	CalleeID       string     `json:"calleeId"`
	StartTime      time.Time  `json:"startTime"`
	EndTime        *time.Time `json:"endTime,omitempty"`
	Status         CallStatus `json:"status"`
}

type SignalingService struct {
	hub      *ws.Hub
	sessions map[string]*CallSession
	mu       sync.RWMutex
}

func NewSignalingService(hub *ws.Hub) *SignalingService {
	service := &SignalingService{
		hub:      hub,
		sessions: make(map[string]*CallSession),
	}

	hub.RegisterHandler(ws.MessageTypeOffer, service.HandleOffer)
	hub.RegisterHandler(ws.MessageTypeAnswer, service.HandleAnswer)
	hub.RegisterHandler(ws.MessageTypeICE, service.HandleICECandidate)
	hub.RegisterHandler(ws.MessageTypeLeave, service.HandleLeave)

	go service.startCleanupScheduler()

	log.Println("📞 Signaling service initialized with WebRTC handlers")
	return service
}

func (s *SignalingService) startCleanupScheduler() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	log.Println("Cleanup scheduler started (runs every 10 minutes)")
	for range ticker.C {
		log.Println("Running session cleanup...")
		s.CleanupOldSessions()
	}
}

func (s *SignalingService) HandleOffer(hub *ws.Hub, message ws.Message) {
	log.Printf(" Handling offer from %s to %s", message.SenderID, message.ReceiverID)

	var sdpData SDPData
	if err := json.Unmarshal(message.Data, &sdpData); err != nil {
		log.Printf(" Error parsing SDP offer: %v", err)
		return
	}

	sessionID := uuid.New().String()
	session := &CallSession{
		SessionID:      sessionID,
		ConversationID: message.ConversationID,
		CallerID:       message.SenderID,
		CalleeID:       message.ReceiverID,
		StartTime:      time.Now(),
		Status:         CallStatusRinging,
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	offerMessage := ws.Message{
		Type:           ws.MessageTypeOffer,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Data:           message.Data,
	}

	dataWithSession := map[string]interface{}{
		"sdp":       sdpData.SDP,
		"type":      sdpData.Type,
		"sessionId": sessionID,
	}
	dataBytes, _ := json.Marshal(dataWithSession)
	offerMessage.Data = dataBytes

	hub.SendToClient(message.ConversationID, message.ReceiverID, offerMessage)
	log.Printf(" Offer forwarded to %s, session: %s", message.ReceiverID, sessionID)
}

func (s *SignalingService) HandleAnswer(hub *ws.Hub, message ws.Message) {
	log.Printf(" Handling answer from %s to %s", message.SenderID, message.ReceiverID)

	var answerData map[string]interface{}
	if err := json.Unmarshal(message.Data, &answerData); err != nil {
		log.Printf(" Error parsing answer: %v", err)
		return
	}

	sessionID, ok := answerData["sessionId"].(string)
	if !ok {
		log.Printf("No session ID in answer")
		return
	}

	s.mu.Lock()
	session, exists := s.sessions[sessionID]
	if !exists {
		s.mu.Unlock()
		log.Printf("Session %s not found", sessionID)
		return
	}

	if session.ConversationID != message.ConversationID {
		s.mu.Unlock()
		log.Printf("Session conversation mismatch")
		return
	}
	if message.SenderID != session.CalleeID {
		s.mu.Unlock()
		log.Printf("Answer sender is not callee")
		return
	}

	session.Status = CallStatusActive
	s.mu.Unlock()
	log.Printf(" Call session %s is now active", sessionID)

	answerMessage := ws.Message{
		Type:           ws.MessageTypeAnswer,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Data:           message.Data,
	}

	hub.SendToClient(message.ConversationID, message.ReceiverID, answerMessage)
	log.Printf(" Answer forwarded to %s", message.ReceiverID)
}

func (s *SignalingService) HandleICECandidate(hub *ws.Hub, message ws.Message) {
	log.Printf("Handling ICE candidate from %s to %s", message.SenderID, message.ReceiverID)

	iceMessage := ws.Message{
		Type:           ws.MessageTypeICE,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Data:           message.Data,
	}

	hub.SendToClient(message.ConversationID, message.ReceiverID, iceMessage)
	log.Printf(" ICE candidate forwarded to %s", message.ReceiverID)
}

func (s *SignalingService) HandleLeave(hub *ws.Hub, message ws.Message) {
	log.Printf(" Handling leave from %s", message.SenderID)

	var leaveData map[string]interface{}
	if err := json.Unmarshal(message.Data, &leaveData); err != nil {
		log.Printf(" Error parsing leave data: %v", err)
		return
	}

	sessionID, ok := leaveData["sessionId"].(string)
	if !ok {
		log.Printf("No session ID in leave message")
	} else {

		s.mu.Lock()
		session, exists := s.sessions[sessionID]
		if !exists {
			s.mu.Unlock()
			log.Printf("Session %s not found", sessionID)
		} else {

			if message.SenderID != session.CallerID && message.SenderID != session.CalleeID {
				s.mu.Unlock()
				log.Printf("Leave sender not authorized")
				return
			}

			now := time.Now()
			session.EndTime = &now
			session.Status = CallStatusEnded
			s.mu.Unlock()
			log.Printf(" Call session %s ended", sessionID)
		}
	}

	leaveMessage := ws.Message{
		Type:           ws.MessageTypeLeave,
		ConversationID: message.ConversationID,
		SenderID:       message.SenderID,
		ReceiverID:     message.ReceiverID,
		Data:           message.Data,
	}

	hub.SendToClient(message.ConversationID, message.ReceiverID, leaveMessage)
	log.Printf(" Leave notification sent to %s", message.ReceiverID)
}

func (s *SignalingService) GetActiveSession(conversationID string) *CallSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, session := range s.sessions {
		if session.ConversationID == conversationID && session.Status == CallStatusActive {
			return session
		}
	}
	return nil
}

func (s *SignalingService) GetSession(sessionID string) *CallSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[sessionID]
}

func (s *SignalingService) CleanupOldSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-1 * time.Hour)
	for id, session := range s.sessions {
		if session.Status == CallStatusEnded && session.EndTime != nil && session.EndTime.Before(cutoff) {
			delete(s.sessions, id)
			log.Printf("Cleaned up old session: %s", id)
		}
	}
}
