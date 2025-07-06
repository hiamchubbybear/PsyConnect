package dto

type SwipeRequest struct {
	TherapistId string `json:"therapist_id" binding:"required"`
}
