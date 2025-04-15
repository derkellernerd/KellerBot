package handler

import "time"

type ErrorResponse struct {
	Error     string
	Timestamp time.Time
}

type SuccessResponse struct {
	Data      any
	Timestamp time.Time
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:     err.Error(),
		Timestamp: time.Now(),
	}
}

func NewSuccessResponse(data any) SuccessResponse {
	return SuccessResponse{
		Data:      data,
		Timestamp: time.Now(),
	}
}
