package apiresponse

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewApiResponse(c *gin.Context, data interface{}) {

	c.JSON(
		http.StatusOK,
		ApiResponse{
			Message: "success",
			Status:  http.StatusOK,
			Data:    data,
		})
}

func ErrorHandler(c *gin.Context, status int, message string) {
	c.JSON(status,
		ApiResponse{
			Message: message,
			Status:  status,
		})
}
