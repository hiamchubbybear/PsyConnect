package http

import (
	"consultationservice/internal/matching/usecase"
	"consultationservice/internal/utils"
	"consultationservice/pkg/apiresponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createMatchUC *usecase.CreateMatchUseCase
	getMatchesUC  *usecase.GetClientMatchesUseCase
	converter     utils.Converter
}

func NewHandler(
	createMatchUC *usecase.CreateMatchUseCase,
	getMatchesUC *usecase.GetClientMatchesUseCase,
) *Handler {
	return &Handler{
		createMatchUC: createMatchUC,
		getMatchesUC:  getMatchesUC,
	}
}

type CreateMatchRequest struct {
	ClientID    string   `json:"client_id"` 
	TherapistID string   `json:"therapist_id" binding:"required"`
	Source      string   `json:"source"`
	SwipeScore  float64  `json:"swipe_score"`
	Reasons     []string `json:"reasons"`
}

func (h *Handler) MatchRequest(c *gin.Context) {
	var req CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid input")
		return
	}

	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		
		
		
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Missing profile ID")
		return
	}

	err := h.createMatchUC.Execute(
		c.Request.Context(),
		profileID,
		req.TherapistID,
		req.Source,
		req.SwipeScore,
		req.Reasons,
	)

	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, "Match success")
}

func (h *Handler) GetAllMatchTherapist(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Missing profile ID")
		return
	}

	pageStr := c.Query("page")
	page, err := h.converter.StringToInt64(pageStr)
	if err != nil {
		page = 1
	}

	matches, err := h.getMatchesUC.Execute(c.Request.Context(), profileID, page)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, "Failed to load match list")
		return
	}

	apiresponse.NewApiResponse(c, matches)
}
