package handlers

import (
	"consultationservice/internal/dto"
	"consultationservice/internal/kafka"
	"consultationservice/internal/repository"
	"consultationservice/pkg/apiresponse"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionRepo *repository.SessionRepository
	kafka       *kafka.Producer
}

func (r *SessionHandler) CreateNewSessionHandler(c *gin.Context) {
	var (
		request dto.SessionRequest
	)
	err := c.ShouldBind(&request)
	if err != nil {
		apiresponse.ErrorHandler(c, 404, "Invalid input")
		return
	}
	res, err := r.sessionRepo.CreateNewSession(request)
	if err != nil {
		r.kafka.SendLogs(err.Error())
		return
	}
	apiresponse.NewApiResponse(c, res)
	return

}
