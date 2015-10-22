package handler

// ApiError defines structure of error that can be returned by the API
type ApiError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewApiError(message string) ApiError {
	return ApiError{Message: message}
}
