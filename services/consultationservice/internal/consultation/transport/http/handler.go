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
	createSessionUC *usecase.CreateSessionUseCase
	getSessionUC    *usecase.GetSessionUseCase
	deleteSessionUC *usecase.DeleteSessionUseCase
	startCallUC     *usecase.StartCallUseCase
}

func NewHandler(
	createSessionUC *usecase.CreateSessionUseCase,
	getSessionUC *usecase.GetSessionUseCase,
	deleteSessionUC *usecase.DeleteSessionUseCase,
	startCallUC *usecase.StartCallUseCase,
) *Handler {
	return &Handler{
		createSessionUC: createSessionUC,
		getSessionUC:    getSessionUC,
		deleteSessionUC: deleteSessionUC,
		startCallUC:     startCallUC,
	}
}

type CreateSessionRequest struct {
	TherapistID string                  `json:"therapist_id" binding:"required"`
	ClientID    string                  `json:"client_id" binding:"required"`
	Mode        domain.ConsultationMode `json:"mode" binding:"required"`
	StartTime   time.Time               `json:"start_time" binding:"required"`
	EndTime     time.Time               `json:"end_time" binding:"required"`
	TimeZone    string                  `json:"time_zone"binding:"required"`
	Price       float64                 `json:"price" binding:"required,gt=0"`
}

// Only client can be use this
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

	ucReq := usecase.CreateSessionRequest{
		TherapistID: req.TherapistID,
		ClientID:    req.ClientID,
		Mode:        req.Mode,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Price:       req.Price,
		TimeZone:    req.TimeZone,
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

// DeleteSession Deprecated
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

//POST /sessions/{id}/payment
//Activate    POST /sessions/{id}/activate
//Start call    POST /sessions/{id}/call/start
//Complete    POST /sessions/{id}/complete
//Cancel    POST /sessions/{id}/cancel
