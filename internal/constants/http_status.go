package constants

import "net/http"

const (
	StatusOK                  = http.StatusOK
	StatusCreated             = http.StatusCreated
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusForbidden           = http.StatusForbidden
	StatusNotFound            = http.StatusNotFound
	StatusConflict            = http.StatusConflict
	StatusUnprocessableEntity = http.StatusUnprocessableEntity
	StatusTooManyRequests     = http.StatusTooManyRequests
	StatusInternalServerError = http.StatusInternalServerError
	StatusNotImplemented      = http.StatusNotImplemented
	StatusServiceUnavailable  = http.StatusServiceUnavailable
)

const (
	MessageSuccess             = "Success."
	MessageBadRequest          = "Bad Request"
	MessageUnauthorized        = "Unauthorized"
	MessageForbidden           = "Forbidden"
	MessageNotFound            = "Resource not found"
	MessageRouteNotFound       = "Route not found"
	MessageConflict            = "Conflict"
	MessageUnprocessableEntity = "Unprocessable Entity"
	MessageInternalServerError = "Internal Server Error"
)
