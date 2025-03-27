package main

import (
	"consultationservice/apiresponse"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewServer() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		response := apiresponse.NewApiResponse("Response from user", http.StatusOK, "Hello World!")
		c.JSON(response.Status, response)
	})
	r.GET("/", func(c *gin.Context) {
		response := apiresponse.NewApiResponse("Response from user", http.StatusOK, "Hello World!")
		c.JSON(response.Status, response)
	})
	r.GET("/", func(c *gin.Context) {
		response := apiresponse.NewApiResponse("Response from user", http.StatusOK, "Hello World!")
		c.JSON(response.Status, response)
	})
	r.GET("/", func(c *gin.Context) {
		response := apiresponse.NewApiResponse("Response from user", http.StatusOK, "Hello World!")
		c.JSON(response.Status, response)
	})
	r.GET("/", func(c *gin.Context) {
		response := apiresponse.NewApiResponse("Response from user", http.StatusOK, "Hello World!")
		c.JSON(response.Status, response)
	})
	r.Run(":8083")
}
