package handlers

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/db"
	"consultationservice/internal/grpc/handler"
	grpc "consultationservice/internal/grpc/handler"
	"consultationservice/internal/model"
	"consultationservice/internal/repository"
	"consultationservice/pkg/apiresponse"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	clientDB    *mongo.Collection
	clientRepo  *repository.ClientRepository
	grpcProfile *grpc.ProfileGrpc
)

type ClientHandler struct {
}

func InitClientHandler(env *bootstrap.Env) {
	clientDB = db.GetClientCollection()
	clientRepo = repository.NewClientRepository(clientDB)
	grpcAddress := env.GrpcAdd
	if grpcAddress == "" {
		log.Fatal("Failed to load env file")
	}
	grpcProfile, _ = handler.NewProfileGrpc(grpcAddress)
}

// External Rest API -- GET /consultation/client
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
	res, err := grpcProfile.CheckProfileExists(profileId)
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
