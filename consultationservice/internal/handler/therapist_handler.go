package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	"consultationservice/internal/grpc/handler"
	"consultationservice/internal/model"
	"consultationservice/internal/repository"
	"consultationservice/pkg/apiresponse"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	therapistDB    *mongo.Collection
	therapistRepo  *repository.TherapistRepository
	grpcRepository *handler.ProfileGrpc
)

type TherapistHandler struct {
}

func InitTherapistHandler(env *bootstrap.Env) {
	therapistDB = db.GetTherapistCollection()
	therapistRepo = repository.NewTherapistRepository(therapistDB)
	grpcAddress := env.GrpcAdd
	if grpcAddress == "" {
		log.Fatal("Failed to load env file")
	}
	grpcRepository, _ = handler.NewProfileGrpc(grpcAddress)
}

// External Rest API -- GET /consultation/{profile_id}
func (r TherapistHandler) GetTherapistHandler(c *gin.Context) {
	profileID := c.Request.Header["X-Profile-Id"][0]
	if profileID == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	res, err := therapistRepo.FindTherapistMatchingProfile(profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)
}

// External Rest API -- POST /consultation/:=profile_id
func (r TherapistHandler) PostTherapistHandler(c *gin.Context) {
	var therapist model.Therapist
	profileId := c.Request.Header["X-Profile-Id"][0]
	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	if err := c.ShouldBindJSON(&therapist); err != nil {
		apiresponse.ErrorHandler(c, 500, "Invalid input")
		return
	}
	if grpcRepository == nil {
		log.Printf("Grpc are not initialize")
	}
	res, err := grpcRepository.CheckProfileExists(profileId)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}
	if !res {
		apiresponse.ErrorHandler(c, 404, "Profile not found")
		return
	}
	therapist.ProfileId = profileId
	_, err = therapistRepo.CreateTherapistMatchingProfile(&therapist)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, therapist)
}

// External Rest API -- POST /consultation
// Header -X Profile-Id = {id}
func (r TherapistHandler) PutTherapistHandler(c *gin.Context) {
	var therapist *model.Therapist
	profileId := c.Request.Header["X-Profile-Id"][0]
	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	if err := c.ShouldBindJSON(&therapist); err != nil {
		apiresponse.ErrorHandler(c, 400, "Invalid input")
		return
	}
	res, err := therapistRepo.UpdateMatchingProfile(profileId, therapist)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)

}

func (r TherapistHandler) ChangeTherapistProfileStatus(c *gin.Context) {
	profileId := c.GetHeader("X-Profile-Id")
	statusStr := c.Param("status")
	status, err := strconv.ParseBool(statusStr)

	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	object, err := therapistRepo.FindTherapistMatchingProfile(profileId)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, err.Error())
		return
	}
	if object.CurrentSession != nil {
		apiresponse.ErrorHandler(c, 404, "You can not disable account please complete all your current schedules ")
		return
	}
	res, err := therapistRepo.DisableTherapistMatchingProfile(profileId, status)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}
	if !res {
		apiresponse.ErrorHandler(c, 404, "Failed to disable account")
	}
	apiresponse.NewApiResponse(c, status)
}
