package dto

import "consultationservice/internal/model"

type FilterRawData struct {
	ClientRaw    model.Client      `json:"clientRaw"`
	TherapistRaw []model.Therapist `json:"therapistsRaw"`
}
type FilterRawDataV1 struct {
	ClientRaw    model.Client        `json:"clientRaw"`
	TherapistRaw []model.TherapistV1 `json:"therapistsRaw"`
}
