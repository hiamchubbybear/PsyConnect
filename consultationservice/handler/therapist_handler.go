package handlers

import (
	"consultationservice/apiresponse"
	"consultationservice/db"
	"consultationservice/grpc/handler"
	"consultationservice/model"
	"consultationservice/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var (
	therapistDB   *mongo.Collection
	therapistRepo *repository.TherapistRepository
)

func InitTherapistHandler() {
	therapistDB = db.GetTherapistCollection()
	therapistRepo = repository.NewTherapistRepository(therapistDB)
}

// External Rest API -- GET /consultation/{profile_id}
func GetTherapistHandler(c *gin.Context) {
	profileID := c.Param("profile_id")
	res, err := therapistRepo.FindTherapistMatchingProfile(profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)
}

// External Rest API -- POST /consultation/:=profile_id
func PostTherapistHandler(c *gin.Context) {
	var therapist model.Therapist
	profile_id := c.Param("profile_id")
	if profile_id == "" {
		apiresponse.ErrorHandler(c, 400, "Missing path variable")
		return
	}
	if err := c.ShouldBindJSON(&therapist); err != nil {
		apiresponse.ErrorHandler(c, 500, "Invalid input")
		return
	}
	res, err := handler.CheckProfileExists(profile_id)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}
	if !res {
		apiresponse.ErrorHandler(c, 404, "Profile not found")
		return
	}
	therapist.ProfileId = profile_id
	_, err = therapistRepo.CreateTherapistMatchingProfile(&therapist)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, therapist)
}
func gRPCProfileChecking(c *gin.Context) {

}
