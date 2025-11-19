package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/model"
	"consultationservice/internal/repository"
	"consultationservice/internal/utils"
	"consultationservice/pkg/apiresponse"

	"github.com/gin-gonic/gin"
)

var converter utils.Converter

type MatchHandler struct {
	RepoManager *repository.RepositoryManager
}

func NewMatchHandler(env *bootstrap.Env, repoManager *repository.RepositoryManager) *MatchHandler {
	return &MatchHandler{
		RepoManager: repoManager,
	}
}

func (r *MatchHandler) GetAllMatchTherapist(c *gin.Context) {
	profileId := c.GetHeader("X-Profile-Id")
	if profileId == "" {
		apiresponse.ErrorHandler(c, 400, "Missing profile ID")
		return
	}

	pageStr := c.Query("page")
	page, err := converter.StringToInt64(pageStr)
	if err != nil {
		apiresponse.ErrorHandler(c, 400, "Invalid page number")
		return
	}
	matches, err := r.RepoManager.MatchingRepo.FilterAllTherapist(profileId, page)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, "Failed to load match list")
		return
	}

	apiresponse.NewApiResponse(c, matches)
}
func (r *MatchHandler) MatchRequest(c *gin.Context) {
	var request model.Match
	err := c.ShouldBindJSON(&request)
	profileId := c.GetHeader("X-Profile-Id")
	if err != nil {
		apiresponse.ErrorHandler(c, 404, "Invalid input")
		return
	}
	_, err = r.RepoManager.TherapistRepo.FindTherapistMatchingProfile(request.TherapistID)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, "Therapist not found "+err.Error())
		return
	}
	_, err = r.RepoManager.ClientRepo.FindClientMatchingProfile(profileId)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, "Client not found "+err.Error())
		return
	}
	err = r.RepoManager.MatchingRepo.CreateMatch(request)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, "Match success")

}

// Deprecated: This func can not be use anymore ,  Reason : Scale -> Change last commit : 4fc7730
// func (r *MatchHandler) ResponseMatchingRequest(c *gin.Context) {
// 	var request dto.ResponseMatchingRequest
// 	err := c.ShouldBindJSON(&request)
// 	if err != nil {
// 		apiresponse.ErrorHandler(c, 404, "Invalid input")
// 		return
// 	}
// 	therapistId := c.GetHeader("X-Profile-Id")
// 	if therapistId == "" {
// 		apiresponse.ErrorHandler(c, 404, "Profile id not found")
// 		return
// 	}
// 	responseStatus, err := r.RepoManager.GrpcProfile.ResponseMatchingRequest(request, therapistId)
// 	if err != nil {
// 		apiresponse.ErrorHandler(c, 504, err.Error())
// 		return
// 	}
// 	addStatus, err := matchingRepo.AddMatchedClient(therapistId, request.ClientId)
// 	if !addStatus {
// 		apiresponse.ErrorHandler(c, 500, "Failed to add client to therapist")
// 		return
// 	}
// 	if err != nil {
// 		apiresponse.ErrorHandler(c, 500, err.Error())
// 		return
// 	}
// 	apiresponse.NewApiResponse(c, responseStatus)
// 	return

// }
