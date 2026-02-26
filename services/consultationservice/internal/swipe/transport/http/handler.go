package http

import (
	"consultationservice/internal/swipe/usecase"
	"consultationservice/pkg/apiresponse"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	insertSwipeUC   *usecase.InsertSwipeUseCase
	swipeAndMatchUC *usecase.SwipeAndMatchUseCase
	recommendUC     *usecase.RecommendUseCase
}

func NewHandler(
	insertSwipeUC *usecase.InsertSwipeUseCase,
	swipeAndMatchUC *usecase.SwipeAndMatchUseCase,
	recommendUC *usecase.RecommendUseCase,
) *Handler {
	return &Handler{
		insertSwipeUC:   insertSwipeUC,
		swipeAndMatchUC: swipeAndMatchUC,
		recommendUC:     recommendUC,
	}
}

type SwipeTherapistRequest struct {
	ClientID    string   `json:"client_id"`
	TherapistID string   `json:"therapist_id" binding:"required"`
	Points      float32  `json:"points"`
	Reasons     []string `json:"reasons"`
}

func (h *Handler) SwipeTherapist(c *gin.Context) {
	var req SwipeTherapistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Print("failed to bind data")
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Use header profile ID if not in body
	profileID := c.GetHeader("X-Profile-Id")
	if profileID != "" {
		req.ClientID = profileID
	}

	if req.ClientID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Missing client ID")
		return
	}

	err := h.insertSwipeUC.Execute(c.Request.Context(), req.ClientID, req.TherapistID, req.Points, req.Reasons)
	if err != nil {
		log.Print("failed to insert swipe:", err)
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, "Failed to insert swipe")
		return
	}

	apiresponse.NewApiResponse(c, true)
}

// TriggerUpdate - DISABLED: Recommendation feature temporarily disabled
func (h *Handler) TriggerUpdate(c *gin.Context) {
	apiresponse.ErrorHandler(c, http.StatusServiceUnavailable, "Recommendation feature temporarily disabled during migration")
}

// TriggerUpdateV1 fetches and triggers matching update
func (h *Handler) TriggerUpdateV1(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Missing profile ID")
		return
	}

	err := h.recommendUC.TriggerUpdateV1(c.Request.Context(), profileID)
	if err != nil {
		log.Print("Failed to trigger update:", err)
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, "Failed to get recommendations")
		return
	}

	apiresponse.NewApiResponse(c, true)
}

// PopTop5 - DISABLED: Recommendation feature temporarily disabled
func (h *Handler) PopTop5(c *gin.Context) {
	apiresponse.ErrorHandler(c, http.StatusServiceUnavailable, "Recommendation feature temporarily disabled during migration")
}

// PopTop5V1 fetches top 5 swipes from the DB
func (h *Handler) PopTop5V1(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Missing profile ID")
		return
	}

	swipes, err := h.recommendUC.PopTop5V1(c.Request.Context(), profileID)
	if err != nil {
		log.Print("Failed to pop top 5:", err)
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, "Failed to get top recommendations")
		return
	}

	apiresponse.NewApiResponse(c, swipes)
}
