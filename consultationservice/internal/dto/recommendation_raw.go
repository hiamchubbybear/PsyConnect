package dto

import "consultationservice/internal/model"

type FilterRawData struct {
	ClientRaw    model.Client      `json:"clienRaw"`
	TherapistRaw []model.Therapist `json:"therapistsRaw"`
}
