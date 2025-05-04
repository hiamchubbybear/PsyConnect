package handlers

import (
	"consultationservice/internal/db"
	"consultationservice/internal/dto"
	"consultationservice/internal/grpc/handler"
	"consultationservice/internal/repository"
	"consultationservice/internal/utils"
	"consultationservice/pkg/apiresponse"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	matchingRepo *repository.MatchingRepository
	converter    utils.Converter
)

type MatchHandler struct{}

func InitMatchHandler() {
	clientRepo := repository.NewClientRepository(db.GetClientCollection())
	grpcProfile, err := handler.NewProfileGrpc("127.0.0.1:9091")
	if err != nil {
		log.Printf("Error while create grpcConnection %v", err)
	}
	matchingRepo = repository.NewMatchingRepository(db.GetTherapistCollection(), clientRepo, grpcProfile)
}

// GET /consultation/therapist/match?page=1
func (r MatchHandler) GetAllMatchTherapist(c *gin.Context) {
	pageStr := c.Query("page")
	page, err := converter.StringToInt64(pageStr)
	if err != nil {
		apiresponse.ErrorHandler(c, 400, "Invalid page number")
		return
	}

	profileId := c.GetHeader("X-Profile-Id")
	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}

	res, err := matchingRepo.FilterAllMatching(profileId, page)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, res)
}
func (r MatchHandler) MatchRequest(c *gin.Context) {
	var request dto.MatchingRequest
	err := c.ShouldBindJSON(&request)
	profileId := c.GetHeader("X-Profile-Id")
	if err != nil {
		apiresponse.ErrorHandler(c, 404, "Invalid input")
		return
	}
	_, err = therapistRepo.FindTherapistMatchingProfile(request.ThearpistId)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, "Therapist not found "+err.Error())
		return
	}
	_, err = clientRepo.FindClientMatchingProfile(profileId)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, "Client not found "+err.Error())
		return
	}
	res, err := matchingRepo.RequestMatching(request, profileId)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)
}
