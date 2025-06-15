package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/model"
	"consultationservice/internal/repository"
	"consultationservice/pkg/apiresponse"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TherapistHandler struct {
	RepoManager *repository.RepositoryManager
}

func NewTherapistHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager) *TherapistHandler {
	return &TherapistHandler{
		RepoManager: repoManager,
	}
}

func (h *TherapistHandler) GetTherapistHandler(c *gin.Context) {
	profileID := c.GetHeader("X-Profile-Id")
	if profileID == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}

	res, err := h.RepoManager.TherapistRepo.FindTherapistMatchingProfile(profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)
}

func (h *TherapistHandler) PostTherapistHandler(c *gin.Context) {
	var therapist model.Therapist
	profileId := c.GetHeader("X-Profile-Id")
	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}

	if err := c.ShouldBindJSON(&therapist); err != nil {
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

	therapist.ProfileId = profileId
	_, err = h.RepoManager.TherapistRepo.CreateTherapistMatchingProfile(&therapist)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, therapist)
}

func (h *TherapistHandler) PutTherapistHandler(c *gin.Context) {
	var therapist *model.Therapist
	profileId := c.GetHeader("X-Profile-Id")
	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}

	if err := c.ShouldBindJSON(&therapist); err != nil {
		apiresponse.ErrorHandler(c, 400, "Invalid input")
		return
	}

	res, err := h.RepoManager.TherapistRepo.UpdateMatchingProfile(profileId, therapist)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, res)
}

func (h *TherapistHandler) ChangeTherapistProfileStatus(c *gin.Context) {
	profileId := c.GetHeader("X-Profile-Id")
	statusStr := c.Param("status")
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		apiresponse.ErrorHandler(c, 400, "Invalid status parameter")
		return
	}

	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}

	object, err := h.RepoManager.TherapistRepo.FindTherapistMatchingProfile(profileId)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, err.Error())
		return
	}

	if object.CurrentSession != nil {
		apiresponse.ErrorHandler(c, 400, "You cannot disable account, please complete all your current schedules")
		return
	}

	res, err := h.RepoManager.TherapistRepo.DisableTherapistMatchingProfile(profileId, status)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}

	if !res {
		apiresponse.ErrorHandler(c, 404, "Failed to disable account")
		return
	}

	apiresponse.NewApiResponse(c, status)
}
