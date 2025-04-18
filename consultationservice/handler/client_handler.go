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
	clientDB   *mongo.Collection
	clientRepo *repository.ClientRepository
)

type ClientHandler struct {
}

func InitClientHandler() {
	clientDB = db.GetClientCollection()
	clientRepo = repository.NewClientRepository(clientDB)
}

// External Rest API -- GET /consultation/client
// Header -X Profile-Id = {id}
func (r ClientHandler) GetClientHandler(c *gin.Context) {
	profileID := c.Request.Header["X-Profile-Id"][0]
	if profileID == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	res, err := clientRepo.FindClientMatchingProfile(profileID)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)
}

// External Rest API -- POST /consultation/client
// Header -X Profile-Id = {id}
func (r ClientHandler) PostClientHandler(c *gin.Context) {
	var client model.Client
	profileId := c.Request.Header["X-Profile-Id"][0]
	if profileId == "" {
		apiresponse.ErrorHandler(c, 404, "Your token is unavailable or profile id not found")
		return
	}
	if err := c.ShouldBindJSON(&client); err != nil {
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
	client.ProfileId = profileId
	_, err = clientRepo.CreateClientMatchingProfile(&client)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, client)
}

// External Rest API -- POST /consultation/client

func (r ClientHandler) PutClientHandler(c *gin.Context) {
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
	res, err := clientRepo.UpdateMatchingProfile(profileId, client)
	if err != nil {
		apiresponse.ErrorHandler(c, 500, err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)

}
