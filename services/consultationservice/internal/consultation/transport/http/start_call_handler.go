package http

import (
	"consultationservice/pkg/apiresponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) StartCall(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		apiresponse.ErrorHandler(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	callerID := c.GetHeader("X-Profile-Id")
	if callerID == "" {
		apiresponse.ErrorHandler(c, http.StatusUnauthorized, "User authentication required")
		return
	}

	if err := h.startCallUC.Execute(c.Request.Context(), sessionID, callerID); err != nil {
		apiresponse.ErrorHandler(c, http.StatusInternalServerError, err.Error())
		return
	}

	apiresponse.NewApiResponse(c, gin.H{
		"message": "Call initiated successfully",
	})
}
