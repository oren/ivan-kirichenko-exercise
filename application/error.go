package application

import "net/http"

// ApiError defines structure of error that can be returned by the API
type ApiError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var ErrorUnautorized ApiError = ApiError{http.StatusUnauthorized, "no or incorrect authorization token provided", nil}

func NewApiError(code int, message string, data interface{}) ApiError {
	return ApiError{code, message, data}
}
