package handlers

import (
	"consultationservice/apiresponse"
	"consultationservice/db"
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
	therapistDB = db.InitTherapistDB()
	therapistRepo = repository.NewTherapistRepository(therapistDB)
}

// GET /therapist
func GetTherapist(c *gin.Context) {
	response := apiresponse.NewApiResponse("Therapist data fetched", http.StatusOK, gin.H{
		"id":   "1",
		"name": "Dr. Alice",
	})
	c.JSON(response.Status, response)
}

// POST /therapist
func PostTherapist(c *gin.Context) {
	var therapist model.Therapist
	if err := c.ShouldBindJSON(&therapist); err != nil {
		response := apiresponse.NewErrorHandler("Invalid input", http.StatusBadRequest)
		c.JSON(response.Status, response)
		return
	}

	result, err := therapistRepo.CreateTherapistMatchingProfile(&therapist)
	if err != nil {
		response := apiresponse.NewErrorHandler(err.Error(), http.StatusInternalServerError)
		c.JSON(response.Status, response)
		return
	}

	response := apiresponse.NewApiResponse("Therapist created successfully", http.StatusCreated, result)
	c.JSON(response.Status, response)
}
