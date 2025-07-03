package dto

import "consultationservice/internal/model"

type FilterRawData struct {
	ClientRaw    model.Client      `json:"clientRaw"`
	TherapistRaw []model.Therapist `json:"therapistsRaw"`
}
