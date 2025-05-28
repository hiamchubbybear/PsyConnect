package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/dto"
	"consultationservice/internal/repository"
	"consultationservice/pkg/apiresponse"
	"log"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	RepoManager *repository.RepositoryManager
}

func NewSessionHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager) *SessionHandler {
	return &SessionHandler{
		RepoManager: repoManager,
	}
}

func (h *SessionHandler) CreateNewSessionHandler(c *gin.Context) {
	var request dto.SessionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid request body: %v", err)
		apiresponse.ErrorHandler(c, 400, "Invalid input format")
		return
	}

	if h.RepoManager == nil {
		log.Printf("RepoManager is nil")
		apiresponse.ErrorHandler(c, 500, "Internal server error")
		return
	}

	if h.RepoManager.SessionRepo == nil {
		log.Printf("SessionRepo is nil")
		apiresponse.ErrorHandler(c, 500, "Session repository not initialized")
		return
	}

	log.Printf("Creating session request: ClientID=%s, TherapistID=%s, SessionTime=%s",
		request.ClientID, request.TherapistID, request.SessionTime)
	result, err := h.RepoManager.SessionRepo.CreateNewSession(request)
	if err != nil {
		log.Printf("Failed to create session: %v", err)

		if h.RepoManager.Kafka != nil {
			h.RepoManager.Kafka.SendLogs(err.Error())
		}

		apiresponse.ErrorHandler(c, 400, err.Error())
		return
	}

	log.Printf("Session created successfully: %+v", result)
	apiresponse.NewApiResponse(c, request)
}
func (h *SessionHandler) DeleteCurrentSessionHandler(c *gin.Context) {
	var request dto.DeleteSessionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid delete request body: %v", err)
		apiresponse.ErrorHandler(c, 400, "Invalid input format")
		return
	}

	if h.RepoManager == nil || h.RepoManager.SessionRepo == nil {
		log.Println("Session repository not initialized")
		apiresponse.ErrorHandler(c, 500, "Internal server error")
		return
	}

	success, err := h.RepoManager.SessionRepo.DeleteCurrentSession(request)
	if err != nil {
		log.Printf("Failed to delete session: %v", err)
		if h.RepoManager.Kafka != nil {
			h.RepoManager.Kafka.SendLogs("Delete session failed: " + err.Error())
		}
		apiresponse.ErrorHandler(c, 400, err.Error())
		return
	}

	if !success {
		log.Println("Session deletion unsuccessful")
		apiresponse.ErrorHandler(c, 400, "Session deletion failed")
		return
	}
	apiresponse.NewApiResponse(c, gin.H{
		"message":    "Session deleted successfully",
		"session_id": request.SessionID,
	})
}
func (h *SessionHandler) GetAllSessions(c *gin.Context) {
	sessions, err := h.RepoManager.SessionRepo.GetAllSessions()
	if err != nil {
		apiresponse.ErrorHandler(c, 500, "Error fetching sessions")
		return
	}
	apiresponse.NewApiResponse(c, sessions)
}

func (h *SessionHandler) GetSessionByID(c *gin.Context) {
	id := c.Param("id")
	session, err := h.RepoManager.SessionRepo.GetSessionByID(id)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, "Session not found")
		return
	}
	apiresponse.NewApiResponse(c, session)
}
