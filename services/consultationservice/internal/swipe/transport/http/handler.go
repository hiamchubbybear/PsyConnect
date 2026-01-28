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
}

func NewHandler(
	insertSwipeUC *usecase.InsertSwipeUseCase,
	swipeAndMatchUC *usecase.SwipeAndMatchUseCase,
) *Handler {
	return &Handler{
		insertSwipeUC:   insertSwipeUC,
		swipeAndMatchUC: swipeAndMatchUC,
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

// TriggerUpdateV1 - DISABLED: Recommendation feature temporarily disabled
func (h *Handler) TriggerUpdateV1(c *gin.Context) {
	apiresponse.ErrorHandler(c, http.StatusServiceUnavailable, "Recommendation feature temporarily disabled during migration")
}

// PopTop5 - DISABLED: Recommendation feature temporarily disabled
func (h *Handler) PopTop5(c *gin.Context) {
	apiresponse.ErrorHandler(c, http.StatusServiceUnavailable, "Recommendation feature temporarily disabled during migration")
}

// PopTop5V1 - DISABLED: Recommendation feature temporarily disabled
func (h *Handler) PopTop5V1(c *gin.Context) {
	apiresponse.ErrorHandler(c, http.StatusServiceUnavailable, "Recommendation feature temporarily disabled during migration")
}
