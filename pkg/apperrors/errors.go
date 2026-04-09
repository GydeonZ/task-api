package apperrors

import "net/http"

type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrNotFound     = &AppError{Code: "NOT_FOUND", Message: "resource not found"}
	ErrUnauthorized = &AppError{Code: "UNAUTHORIZED", Message: "unauthorized"}
	ErrConflict     = &AppError{Code: "CONFLICT", Message: "resource already exists"}
	ErrBadRequest   = &AppError{Code: "BAD_REQUEST", Message: "bad request"}
)

func HTTPStatus(err error) int {
	appErr, ok := err.(*AppError)
	if !ok {
		return http.StatusInternalServerError
	}
	switch appErr.Code {
	case "NOT_FOUND":
		return http.StatusNotFound
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "CONFLICT":
		return http.StatusConflict
	case "BAD_REQUEST":
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
