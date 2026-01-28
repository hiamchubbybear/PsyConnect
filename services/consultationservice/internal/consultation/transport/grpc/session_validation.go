package grpc

import (
	"consultationservice/internal/consultation/usecase"
	"context"
)

// SessionValidationHandler handles gRPC requests for session validation
// This will be used by Chat Service to validate sessions before starting WebRTC calls
type SessionValidationHandler struct {
	getSessionUC *usecase.GetSessionUseCase
}

// NewSessionValidationHandler creates a new gRPC handler
func NewSessionValidationHandler(getSessionUC *usecase.GetSessionUseCase) *SessionValidationHandler {
	return &SessionValidationHandler{
		getSessionUC: getSessionUC,
	}
}

// ValidateSessionRequest represents the gRPC request
type ValidateSessionRequest struct {
	SessionID string
}

// ValidateSessionResponse represents the gRPC response
type ValidateSessionResponse struct {
	Valid        bool
	Error        string
	SessionID    string
	TherapistID  string
	ClientID     string
	StartTime    int64
	EndTime      int64
	CanStartCall bool
}

// ValidateSession validates if a session exists and is ready for a call
// This method will be called by Chat Service via gRPC
func (h *SessionValidationHandler) ValidateSession(ctx context.Context, req ValidateSessionRequest) (*ValidateSessionResponse, error) {
	// Get session
	session, err := h.getSessionUC.Execute(ctx, req.SessionID)
	if err != nil {
		return &ValidateSessionResponse{
			Valid: false,
			Error: "session not found: " + err.Error(),
		}, nil
	}

	// Check if session is active
	if !session.IsActive() {
		return &ValidateSessionResponse{
			Valid: false,
			Error: "session is not active",
		}, nil
	}

	// Check if call can be started
	canStartCall := session.CanStartCall()

	return &ValidateSessionResponse{
		Valid:        true,
		SessionID:    session.SessionID,
		TherapistID:  session.TherapistID,
		ClientID:     session.ClientID,
		StartTime:    session.StartTime.Unix(),
		EndTime:      session.EndTime.Unix(),
		CanStartCall: canStartCall,
	}, nil
}

// Note: This is a placeholder implementation
// You'll need to generate proper gRPC code from proto definition
// and implement the generated interface
