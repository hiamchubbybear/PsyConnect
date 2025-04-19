package handlers

import (
	"consultationservice/internal/db"
	"consultationservice/internal/grpc/handler"
	"consultationservice/internal/model"
	"consultationservice/internal/repository"
	"consultationservice/pkg/apiresponse"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

var (
	therapistDB   *mongo.Collection
	therapistRepo *repository.TherapistRepository
)

type TherapistHandler struct {
}

func InitTherapistHandler() {
	therapistDB = db.GetTherapistCollection()
	therapistRepo = repository.NewTherapistRepository(therapistDB)
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
	res, err := handler.CheckProfileExists(profileId)
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
