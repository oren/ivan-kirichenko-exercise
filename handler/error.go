package handler

import (
	"encoding/json"
)

// ApiError defines structure of error that can be returned by the API
type ApiError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// String returns string representation of an API error
func (e ApiError) String() string {
	content, _ := json.Marshal(e)
	return string(content)
}

// NewApiError creates new instance of API error
func NewApiError(message string) ApiError {
	return ApiError{Message: message}
}
