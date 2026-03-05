package http

import (
	"consultationservice/internal/consultation/domain"
	"consultationservice/internal/consultation/usecase"
	"consultationservice/pkg/apiresponse"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createSessionUC         *usecase.CreateSessionUseCase
	getSessionUC            *usecase.GetSessionUseCase
	deleteSessionUC         *usecase.DeleteSessionUseCase
	startCallUC             *usecase.StartCallUseCase
	getPaymentUrlUC         *usecase.GetPaymentURLUseCase
	processPaymentWebhookUC *usecase.ProcessPaymentWebhookUseCase
	processRefundUC         *usecase.ProcessRefundUseCase
	getCalendarUC           *usecase.GetCalendarUseCase
	getOverviewUC           *usecase.GetOverviewUseCase
}

func NewHandler(
	createSessionUC *usecase.CreateSessionUseCase,
	getSessionUC *usecase.GetSessionUseCase,
	deleteSessionUC *usecase.DeleteSessionUseCase,
	startCallUC *usecase.StartCallUseCase,
	getPaymentUrlUC *usecase.GetPaymentURLUseCase,
	processPaymentWebhookUC *usecase.ProcessPaymentWebhookUseCase,
	processRefundUC *usecase.ProcessRefundUseCase,
	getCalendarUC *usecase.GetCalendarUseCase,
	getOverviewUC *usecase.GetOverviewUseCase,
) *Handler {
	return &Handler{
		createSessionUC:         createSessionUC,
		getSessionUC:            getSessionUC,
		deleteSessionUC:         deleteSessionUC,
		startCallUC:             startCallUC,
		getPaymentUrlUC:         getPaymentUrlUC,
		processPaymentWebhookUC: processPaymentWebhookUC,
		processRefundUC:         processRefundUC,
		getCalendarUC:           getCalendarUC,
		getOverviewUC:           getOverviewUC,
	}
}

type CreateSessionRequest struct {
	TherapistID   string                  `json:"therapist_id" binding:"required"`
	ClientID      string                  `json:"client_id"`
	Mode          domain.ConsultationMode `json:"mode" binding:"required"`
	StartTime     time.Time               `json:"start_time" binding:"required"`
	EndTime       time.Time               `json:"end_time" binding:"required"`
	TimeZone      string                  `json:"time_zone" binding:"required"`
	ScheduledDate string                  `json:"scheduled_date,omitempty"`
	Price         float64                 `json:"price" binding:"required,gt=0"`
	LocationInfo  *domain.LocationInfo    `json:"location_info,omitempty"`
}


func (h *Handler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	userId := c.GetHeader("X-Profile-Id")
	if userId == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Unauthenticated")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	req.ClientID = userId

	ucReq := usecase.CreateSessionRequest{
		TherapistID:   req.TherapistID,
		ClientID:      req.ClientID,
		Mode:          req.Mode,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Price:         req.Price,
		TimeZone:      req.TimeZone,
		ScheduledDate: req.ScheduledDate,
		LocationInfo:  req.LocationInfo,
	}

	session, err := h.createSessionUC.Execute(c.Request.Context(), ucReq)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, session)
}

func (h *Handler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	session, err := h.getSessionUC.Execute(c.Request.Context(), sessionID)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusNotFound, "Session not found")
		return
	}

	apiresponse.NewApiResponse(c, session)
}

func (h *Handler) GetAllSessions(c *gin.Context) {
	sessions, err := h.getSessionUC.ExecuteAll(c.Request.Context())
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, "Failed to retrieve sessions")
		return
	}

	apiresponse.NewApiResponse(c, sessions)
}

func (h *Handler) GetSessionsByProfile(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	role := c.GetString("roles")

	var sessions interface{}
	var err error

	if role == "role.therapist" {
		sessions, err = h.getSessionUC.ExecuteByTherapist(c.Request.Context(), profileID)
	} else {
		sessions, err = h.getSessionUC.ExecuteByClient(c.Request.Context(), profileID)
	}

	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, "Failed to retrieve sessions")
		return
	}

	apiresponse.NewApiResponse(c, sessions)
}


func (h *Handler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	if err := h.deleteSessionUC.Execute(c.Request.Context(), sessionID); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, gin.H{
		"message":    "Session deleted successfully",
		"session_id": sessionID,
	})
}

func (h *Handler) GetPaymentURL(c *gin.Context) {
	sessionID := c.Param("id")
	url, err := h.getPaymentUrlUC.Execute(c.Request.Context(), sessionID)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, gin.H{"payment_url": url})
}

func (h *Handler) ProcessPaymentWebhook(c *gin.Context) {
	var payload usecase.PaymentWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid payload")
		return
	}
	if err := h.processPaymentWebhookUC.Execute(c.Request.Context(), payload); err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, gin.H{"message": "Webhook processed successfully"})
}

func (h *Handler) RefundSession(c *gin.Context) {
	sessionID := c.Param("id")
	traceID, err := h.processRefundUC.Execute(c.Request.Context(), sessionID)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, gin.H{"message": "Refund successful", "refund_trace_id": traceID})
}


func (h *Handler) GetCalendar(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	fromStr := c.DefaultQuery("from", "")
	toStr := c.DefaultQuery("to", "")
	role := c.GetString("roles")

	events, err := h.getCalendarUC.Execute(c.Request.Context(), profileID, role, fromStr, toStr)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, "Failed to retrieve calendar")
		return
	}

	apiresponse.NewApiResponse(c, gin.H{"events": events})
}


func (h *Handler) GetOverview(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "Profile ID is required")
		return
	}

	role := c.GetString("roles")
	overview, err := h.getOverviewUC.Execute(c.Request.Context(), profileID, role)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, "Failed to load overview")
		return
	}

	apiresponse.NewApiResponse(c, overview)
}
