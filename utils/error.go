package utils

import "net/http"

func BadRequest(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

func NotFound(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

func Internal(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}
func Forbidden(message string, err error) *AppError {
	return NewAppError(http.StatusForbidden, message, err)
}
func Unauthorized(message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, message, err)
}