package helpers

import (
	"errors"

	"go-gin-boilerplate/internal/constants"
)

// AppError merepresentasikan operational error di aplikasi (ekuivalen dengan AppError di POS-nextjs-express).
type AppError struct {
	Name          string      `json:"-"`
	Message       string      `json:"message"`
	StatusCode    int         `json:"-"`
	IsOperational bool        `json:"-"`
	Details       interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(name, message string, statusCode int, details ...interface{}) *AppError {
	var det interface{}
	if len(details) > 0 {
		det = details[0]
	}

	return &AppError{
		Name:          name,
		Message:       message,
		StatusCode:    statusCode,
		IsOperational: true,
		Details:       det,
	}
}

// ============================================
// 4xx — Client Errors
// ============================================

func NewErrorBadRequest(message string, details ...interface{}) *AppError {
	if message == "" {
		message = constants.MessageBadRequest
	}
	return NewAppError("ErrorBadRequest", message, constants.StatusBadRequest, details...)
}

func NewErrorUnauthorized(message string, details ...interface{}) *AppError {
	if message == "" {
		message = constants.MessageUnauthorized
	}
	return NewAppError("ErrorUnauthorized", message, constants.StatusUnauthorized, details...)
}

func NewErrorForbidden(message string, details ...interface{}) *AppError {
	if message == "" {
		message = constants.MessageForbidden
	}
	return NewAppError("ErrorForbidden", message, constants.StatusForbidden, details...)
}

func NewErrorNotFound(message string, details ...interface{}) *AppError {
	if message == "" {
		message = constants.MessageNotFound
	}
	return NewAppError("ErrorNotFound", message, constants.StatusNotFound, details...)
}

func NewErrorConflict(message string, details ...interface{}) *AppError {
	if message == "" {
		message = constants.MessageConflict
	}
	return NewAppError("ErrorConflict", message, constants.StatusConflict, details...)
}

func NewErrorValidation(message string, details ...interface{}) *AppError {
	if message == "" {
		message = constants.MessageUnprocessableEntity
	}
	return NewAppError("ErrorValidation", message, constants.StatusUnprocessableEntity, details...)
}

// ============================================
// 5xx — Server Errors
// ============================================

func NewErrorInternalServer(message string, details ...interface{}) *AppError {
	if message == "" {
		message = constants.MessageInternalServerError
	}
	return NewAppError("ErrorInternalServer", message, constants.StatusInternalServerError, details...)
}

// IsAppError memeriksa apakah error merupakan instance dari *AppError.
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
