package http

import (
	"consultationservice/internal/client/domain"
	"consultationservice/internal/client/usecase"
	"consultationservice/pkg/apiresponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createClientUC *usecase.CreateClientUseCase
	getClientUC    *usecase.GetClientUseCase
	updateClientUC *usecase.UpdateClientUseCase
	deleteClientUC *usecase.DeleteClientUseCase
}

func NewHandler(
	createClientUC *usecase.CreateClientUseCase,
	getClientUC *usecase.GetClientUseCase,
	updateClientUC *usecase.UpdateClientUseCase,
	deleteClientUC *usecase.DeleteClientUseCase,
) *Handler {
	return &Handler{
		createClientUC: createClientUC,
		getClientUC:    getClientUC,
		updateClientUC: updateClientUC,
		deleteClientUC: deleteClientUC,
	}
}

type CreateClientRequest struct {
	ProfileID         string   `json:"profile_id"`
	Address           string   `json:"address" binding:"required"`
	Languages         []string `json:"languages" binding:"required"`
	IssueDetail       []string `json:"issue_detail"`
	ConsultationModes []string `json:"consultation_modes" binding:"required"`
	RangePrice        int      `json:"rage_price" binding:"required,gt=0"`
	Availability      struct {
		Days      []string `json:"days" binding:"required"`
		TimeSlots []string `json:"time_slots" binding:"required"`
	} `json:"availability" binding:"required"`
	PreferredTherapistGender   string   `json:"preferred_therapist_gender"`
	ExperienceLevel            string   `json:"experience_level"`
	TherapistSpecialization    []string `json:"specialization"`
	UrgencyLevel               string   `json:"urgency_level"`
	SessionDuration            int      `json:"session_duration"`
	PreferredTherapistLanguage []string `json:"preferred_therapist_language"`
	IsFlexibleWithSchedule     bool     `json:"is_flexible_with_schedule"`
}

func (h *Handler) CreateClient(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required in header")
		return
	}
	req.ProfileID = profileID

	ucReq := usecase.CreateClientRequest{
		ProfileID:                  req.ProfileID,
		Address:                    req.Address,
		Languages:                  req.Languages,
		IssueDetail:                req.IssueDetail,
		ConsultationModes:          req.ConsultationModes,
		RangePrice:                 req.RangePrice,
		AvailabilityDays:           req.Availability.Days,
		AvailabilityTimeSlots:      req.Availability.TimeSlots,
		PreferredTherapistGender:   req.PreferredTherapistGender,
		ExperienceLevel:            req.ExperienceLevel,
		TherapistSpecialization:    req.TherapistSpecialization,
		UrgencyLevel:               req.UrgencyLevel,
		SessionDuration:            req.SessionDuration,
		PreferredTherapistLanguage: req.PreferredTherapistLanguage,
		IsFlexibleWithSchedule:     req.IsFlexibleWithSchedule,
	}

	client, err := h.createClientUC.Execute(c.Request.Context(), ucReq)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, client)
}

func (h *Handler) GetClient(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	client, err := h.getClientUC.Execute(c.Request.Context(), profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusNotFound, "Client not found")
		return
	}

	apiresponse.NewApiResponse(c, client)
}

func (h *Handler) GetAllClients(c *gin.Context) {
	clients, err := h.getClientUC.ExecuteAll(c.Request.Context())
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, "Failed to retrieve clients")
		return
	}

	var response []map[string]interface{}
	for _, client := range clients {
		response = append(response, domainToResponse(client))
	}

	apiresponse.NewApiResponse(c, response)
}

func (h *Handler) GetClientByID(c *gin.Context) {
	profileID := c.Param("id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Profile ID is required")
		return
	}

	client, err := h.getClientUC.Execute(c.Request.Context(), profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusNotFound, "Client not found")
		return
	}

	apiresponse.NewApiResponse(c, client)
}

type UpdateClientRequest struct {
	Address           string   `json:"address"`
	Languages         []string `json:"languages"`
	IssueDetail       []string `json:"issue_detail"`
	ConsultationModes []string `json:"consultation_modes"`
	RangePrice        int      `json:"rage_price"`
	Availability      *struct {
		Days      []string `json:"days"`
		TimeSlots []string `json:"time_slots"`
	} `json:"availability"`
	PreferredTherapistGender   string   `json:"preferred_therapist_gender"`
	ExperienceLevel            string   `json:"experience_level"`
	TherapistSpecialization    []string `json:"specialization"`
	UrgencyLevel               string   `json:"urgency_level"`
	SessionDuration            int      `json:"session_duration"`
	PreferredTherapistLanguage []string `json:"preferred_therapist_language"`
	IsFlexibleWithSchedule     bool     `json:"is_flexible_with_schedule"`
}

func (h *Handler) UpdateClient(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ucReq := usecase.UpdateClientRequest{
		Address:                    req.Address,
		Languages:                  req.Languages,
		IssueDetail:                req.IssueDetail,
		ConsultationModes:          req.ConsultationModes,
		RangePrice:                 req.RangePrice,
		PreferredTherapistGender:   req.PreferredTherapistGender,
		ExperienceLevel:            req.ExperienceLevel,
		TherapistSpecialization:    req.TherapistSpecialization,
		UrgencyLevel:               req.UrgencyLevel,
		SessionDuration:            req.SessionDuration,
		PreferredTherapistLanguage: req.PreferredTherapistLanguage,
		IsFlexibleWithSchedule:     req.IsFlexibleWithSchedule,
	}

	if req.Availability != nil {
		ucReq.AvailabilityDays = req.Availability.Days
		ucReq.AvailabilityTimeSlots = req.Availability.TimeSlots
	}

	if err := h.updateClientUC.Execute(c.Request.Context(), profileID, ucReq); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, gin.H{"message": "Client updated successfully"})
}

func (h *Handler) DeleteClient(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	if err := h.deleteClientUC.Execute(c.Request.Context(), profileID); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, gin.H{"message": "Client deleted successfully"})
}


func domainToResponse(client *domain.Client) map[string]interface{} {
	return map[string]interface{}{
		"profile_id":         client.ProfileID,
		"address":            client.Address,
		"languages":          client.Languages,
		"issue_detail":       client.IssueDetail,
		"consultation_modes": client.ConsultationModes,
		"rage_price":         client.RangePrice,
		"availability": map[string]interface{}{
			"days":       client.Availability.Days,
			"time_slots": client.Availability.TimeSlots,
		},
		"preferred_therapist_gender":   client.PreferredTherapistGender,
		"experience_level":             client.ExperienceLevel,
		"therapist_specialization":     client.TherapistSpecialization,
		"urgency_level":                client.UrgencyLevel,
		"session_duration":             client.SessionDuration,
		"preferred_therapist_language": client.PreferredTherapistLanguage,
		"is_flexible_with_schedule":    client.IsFlexibleWithSchedule,
	}
}
