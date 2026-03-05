package grpc

import (
	"consultationservice/internal/consultation/usecase"
	"context"
)



type SessionValidationHandler struct {
	getSessionUC *usecase.GetSessionUseCase
}


func NewSessionValidationHandler(getSessionUC *usecase.GetSessionUseCase) *SessionValidationHandler {
	return &SessionValidationHandler{
		getSessionUC: getSessionUC,
	}
}


type ValidateSessionRequest struct {
	SessionID string
}


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



func (h *SessionValidationHandler) ValidateSession(ctx context.Context, req ValidateSessionRequest) (*ValidateSessionResponse, error) {
	
	session, err := h.getSessionUC.Execute(ctx, req.SessionID)
	if err != nil {
		return &ValidateSessionResponse{
			Valid: false,
			Error: "session not found: " + err.Error(),
		}, nil
	}

	
	if !session.IsActive() {
		return &ValidateSessionResponse{
			Valid: false,
			Error: "session is not active",
		}, nil
	}

	
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




