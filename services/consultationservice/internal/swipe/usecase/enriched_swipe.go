package usecase

import (
	"consultationservice/internal/model"
)

type EnrichedSwipe struct {
	model.TherapistV1 `json:",inline"`
	Points            float32  `json:"points"`
	Reasons           []string `json:"reasons"`
}
