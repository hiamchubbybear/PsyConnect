package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/model"
	"consultationservice/internal/repository"
	"consultationservice/pkg/apiresponse"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	RepoManager *repository.RepositoryManager
}

func NewClientHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager) *ClientHandler {
	return &ClientHandler{
		RepoManager: repoManager,
	}
}

func (h *ClientHandler) GetClientHandler(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	res, err := h.RepoManager.ClientRepo.FindClientMatchingProfile(profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)
}

func (h *ClientHandler) PostClientHandler(c *gin.Context) {
	var client model.Client
	profileId := c.GetHeader("X-Profile-Id")
	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	if err := c.ShouldBindJSON(&client); err != nil {
		apiresponse.ErrorHandler(c, 400, "Invalid input")
		return
	}

	res, err := h.RepoManager.GrpcProfile.CheckProfileExists(profileId)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}
	if !res {
		apiresponse.ErrorHandler(c, 404, "Profile not found")
		return
	}
	client.ProfileId = profileId
	_, err = h.RepoManager.ClientRepo.CreateClientMatchingProfile(&client)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, client)
}

func (h *ClientHandler) PutClientHandler(c *gin.Context) {
	var client *model.Client
	profileId := c.GetHeader("X-Profile-Id")
	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	if err := c.ShouldBindJSON(&client); err != nil {
		apiresponse.ErrorHandler(c, 400, "Invalid input")
		return
	}
	res, err := h.RepoManager.ClientRepo.UpdateMatchingProfile(profileId, client)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}
	_, error := h.RepoManager.SwipeRepo.FilterAllTherapist(profileId)
	if error != nil {
		log.Println(err)
		// apiresponse.ErrorHandler(c, 400, "Failed to update therapist recommendataion field")
		// return
	}
	apiresponse.NewApiResponse(c, res)
}
func (h *ClientHandler) GetClientByIdHandler(c *gin.Context) {
	profileID := c.Param("id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	res, err := h.RepoManager.ClientRepo.FindClientMatchingProfile(profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)
}
