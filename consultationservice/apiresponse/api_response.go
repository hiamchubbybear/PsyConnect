package apiresponse

type ApiResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

func NewApiResponse(message string, status int, data interface{}) ApiResponse {
	return ApiResponse{
		Message: message,
		Status:  status,
		Data:    data,
	}
}
func NewErrorHandler(message string, status int) ApiResponse {
	return ApiResponse{
		Message: message,
		Status:  status,
	}
}
