package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError is a domain-level error that carries an HTTP status code and a
// human-readable message that is safe to expose to clients.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

// Predefined constructors for common error scenarios.

func New(code int, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, Err: cause}
}

func BadRequest(msg string, cause error) *AppError {
	return New(http.StatusBadRequest, msg, cause)
}

func Unauthorized(msg string) *AppError {
	return New(http.StatusUnauthorized, msg, nil)
}

func Forbidden(msg string) *AppError {
	return New(http.StatusForbidden, msg, nil)
}

func NotFound(resource string) *AppError {
	return New(http.StatusNotFound, fmt.Sprintf("%s not found", resource), nil)
}

func Conflict(msg string) *AppError {
	return New(http.StatusConflict, msg, nil)
}

func Internal(cause error) *AppError {
	return New(http.StatusInternalServerError, "internal server error", cause)
}

// Is allows errors.Is to unwrap AppError.
func Is(err error) (*AppError, bool) {
	var ae *AppError
	ok := errors.As(err, &ae)
	return ae, ok
}
