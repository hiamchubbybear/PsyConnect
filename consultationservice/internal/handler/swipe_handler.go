package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/repository"
	"consultationservice/pkg/apiresponse"
	"log"

	"github.com/gin-gonic/gin"
)

type SwipeHandler struct {
	RepoManager *repository.RepositoryManager
}

func NewSwipeHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager) *SwipeHandler {
	return &SwipeHandler{
		RepoManager: repoManager,
	}
}

func (r *SwipeHandler) TriggerUpdate(c *gin.Context) {
	clientId := c.GetHeader("X-Profile-Id")
	if clientId == "" {
		log.Printf("Client id is nil")
		apiresponse.ErrorHandler(c, 404, "Invalid input")
		return
	}
	res, err := r.RepoManager.SwipeRepo.FilterAllTherapist(clientId)
	if err != nil {
		log.Println(err)
		apiresponse.ErrorHandler(c, 404, "Failed to filter all therapist")
		return
	}
	apiresponse.NewApiResponse(c, res)
}
