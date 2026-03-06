package http

import (
	"consultationservice/internal/location/service"
	"consultationservice/pkg/apiresponse"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	geocodingService *service.GeocodingService
}

func NewHandler(geocodingService *service.GeocodingService) *Handler {
	return &Handler{geocodingService: geocodingService}
}

func (h *Handler) Search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if len(query) < 3 {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Query must be at least 3 characters")
		return
	}

	results, err := h.geocodingService.Search(c.Request.Context(), query)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadGateway, "Failed to search locations")
		return
	}

	apiresponse.NewApiResponse(c, results)
}

func (h *Handler) Reverse(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid latitude")
		return
	}

	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Invalid longitude")
		return
	}

	result, err := h.geocodingService.Reverse(c.Request.Context(), lat, lon)
	if err != nil {
		apiresponse.ErrorHandler(c, http.StatusBadGateway, "Failed to reverse geocode location")
		return
	}

	apiresponse.NewApiResponse(c, result)
}
