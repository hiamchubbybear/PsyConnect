package dto

type CancelSessionRequest struct {
	SessionID string `json:"session_id"`
	Reason    string `json:"reason"`
}
