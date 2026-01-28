package http

import (
	"consultationservice/internal/therapist/usecase"
	"consultationservice/pkg/apiresponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createTherapistUC *usecase.CreateTherapistUseCase
	getTherapistUC    *usecase.GetTherapistUseCase
	updateTherapistUC *usecase.UpdateTherapistUseCase
	deleteTherapistUC *usecase.DeleteTherapistUseCase
}

func NewHandler(
	createTherapistUC *usecase.CreateTherapistUseCase,
	getTherapistUC *usecase.GetTherapistUseCase,
	updateTherapistUC *usecase.UpdateTherapistUseCase,
	deleteTherapistUC *usecase.DeleteTherapistUseCase,
) *Handler {
	return &Handler{
		createTherapistUC: createTherapistUC,
		getTherapistUC:    getTherapistUC,
		updateTherapistUC: updateTherapistUC,
		deleteTherapistUC: deleteTherapistUC,
	}
}

// DTOs
type ProfessionalTitleDTO struct {
	Code    string `json:"code"`
	Display string `json:"display"`
}

type DegreeDTO struct {
	Type        string `json:"type"`
	Field       string `json:"field"`
	Institution string `json:"institution"`
	Year        int    `json:"year"`
}

type CertificationDTO struct {
	Name   string `json:"name"`
	Issuer string `json:"issuer"`
	Year   int    `json:"year"`
}

type ProfessionalInfoDTO struct {
	Title           ProfessionalTitleDTO `json:"title"`
	Degrees         []DegreeDTO          `json:"degrees"`
	Certifications  []CertificationDTO   `json:"certifications"`
	ExperienceYears int                  `json:"experience_years"`
}

type CreateTherapistRequest struct {
	ProfileID         string   `json:"profile_id" binding:"required"`
	Address           string   `json:"address" binding:"required"`
	Languages         []string `json:"languages" binding:"required"`
	Specialization    []string `json:"specialization" binding:"required"`
	ConsultationModes []string `json:"consultation_modes"`
	RagePrice         int      `json:"rage_price" binding:"required,gt=0"`
	Currency          string   `json:"currency"`
	Availability      struct {
		Days      []string `json:"days" binding:"required"`
		TimeSlots []string `json:"time_slots" binding:"required"`
	} `json:"availability" binding:"required"`

	Experience       int                 `json:"experience"`
	Rating           float64             `json:"rating"`
	AvatarOverride   string              `json:"avatar_override"`
	Name             string              `json:"name"`
	ProfessionalInfo ProfessionalInfoDTO `json:"professional_info"`
}

func (h *Handler) CreateTherapist(c *gin.Context) {
	var req CreateTherapistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ucReq := usecase.CreateTherapistRequest{
		ProfileID:             req.ProfileID,
		Address:               req.Address,
		Languages:             req.Languages,
		Specialization:        req.Specialization,
		ConsultationModes:     req.ConsultationModes,
		RagePrice:             req.RagePrice,
		Currency:              req.Currency,
		AvailabilityDays:      req.Availability.Days,
		AvailabilityTimeSlots: req.Availability.TimeSlots,
		Experience:            req.Experience,
		Rating:                req.Rating,
		AvatarOverride:        req.AvatarOverride,
		Name:                  req.Name,
		ProfessionalInfo: usecase.ProfessionalInfo{
			Title: usecase.ProfessionalTitle{
				Code:    req.ProfessionalInfo.Title.Code,
				Display: req.ProfessionalInfo.Title.Display,
			},
			ExperienceYears: req.ProfessionalInfo.ExperienceYears,
		},
	}

	for _, d := range req.ProfessionalInfo.Degrees {
		ucReq.ProfessionalInfo.Degrees = append(ucReq.ProfessionalInfo.Degrees, usecase.Degree{
			Type:        d.Type,
			Field:       d.Field,
			Institution: d.Institution,
			Year:        d.Year,
		})
	}

	for _, cert := range req.ProfessionalInfo.Certifications {
		ucReq.ProfessionalInfo.Certifications = append(ucReq.ProfessionalInfo.Certifications, usecase.Certification{
			Name:   cert.Name,
			Issuer: cert.Issuer,
			Year:   cert.Year,
		})
	}

	therapist, err := h.createTherapistUC.Execute(c.Request.Context(), ucReq)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, therapist)
}

func (h *Handler) GetTherapist(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	therapist, err := h.getTherapistUC.Execute(c.Request.Context(), profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusNotFound, "Therapist not found")
		return
	}

	apiresponse.NewApiResponse(c, therapist)
}

func (h *Handler) GetTherapistByID(c *gin.Context) {
	profileID := c.Param("id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Profile ID is required")
		return
	}

	therapist, err := h.getTherapistUC.Execute(c.Request.Context(), profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusNotFound, "Therapist not found")
		return
	}

	apiresponse.NewApiResponse(c, therapist)
}

type UpdateTherapistRequest struct {
	Address           string   `json:"address"`
	Languages         []string `json:"languages"`
	Specialization    []string `json:"specialization"`
	ConsultationModes []string `json:"consultation_modes"`
	RagePrice         int      `json:"rage_price"`
	Currency          string   `json:"currency"`
	Availability      *struct {
		Days      []string `json:"days"`
		TimeSlots []string `json:"time_slots"`
	} `json:"availability"`
	Experience       int                  `json:"experience"`
	AvatarOverride   string               `json:"avatar_override"`
	Name             string               `json:"name"`
	ProfessionalInfo *ProfessionalInfoDTO `json:"professional_info"`
}

func (h *Handler) UpdateTherapist(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	var req UpdateTherapistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	ucReq := usecase.UpdateTherapistRequest{
		Address:           req.Address,
		Languages:         req.Languages,
		Specialization:    req.Specialization,
		ConsultationModes: req.ConsultationModes,
		RagePrice:         req.RagePrice,
		Currency:          req.Currency,
		Experience:        req.Experience,
		AvatarOverride:    req.AvatarOverride,
		Name:              req.Name,
	}

	if req.Availability != nil {
		ucReq.AvailabilityDays = req.Availability.Days
		ucReq.AvailabilityTimeSlots = req.Availability.TimeSlots
	}

	if req.ProfessionalInfo != nil {
		ucReq.ProfessionalInfo = &usecase.ProfessionalInfo{
			Title: usecase.ProfessionalTitle{
				Code:    req.ProfessionalInfo.Title.Code,
				Display: req.ProfessionalInfo.Title.Display,
			},
			ExperienceYears: req.ProfessionalInfo.ExperienceYears,
		}

		for _, d := range req.ProfessionalInfo.Degrees {
			ucReq.ProfessionalInfo.Degrees = append(ucReq.ProfessionalInfo.Degrees, usecase.Degree{
				Type:        d.Type,
				Field:       d.Field,
				Institution: d.Institution,
				Year:        d.Year,
			})
		}

		for _, cert := range req.ProfessionalInfo.Certifications {
			ucReq.ProfessionalInfo.Certifications = append(ucReq.ProfessionalInfo.Certifications, usecase.Certification{
				Name:   cert.Name,
				Issuer: cert.Issuer,
				Year:   cert.Year,
			})
		}
	}

	if err := h.updateTherapistUC.Execute(c.Request.Context(), profileID, ucReq); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, gin.H{"message": "Therapist updated successfully"})
}

func (h *Handler) UpdateAvailability(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	var req struct {
		IsAvailable bool `json:"is_available"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := h.updateTherapistUC.ExecuteAvailability(c.Request.Context(), profileID, req.IsAvailable); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, gin.H{"message": "Availability updated successfully"})
}

func (h *Handler) DeleteTherapist(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	// Warning: This is a hard delete
	if err := h.deleteTherapistUC.Execute(c.Request.Context(), profileID); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, gin.H{"message": "Therapist deleted successfully"})
}
