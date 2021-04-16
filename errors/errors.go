package errors

import (
	"github.com/blackdreamers/core/config"
	"github.com/blackdreamers/go-micro/v3/errors"
)

// New generates a custom error.
func New(detail string, code int32) error {
	return errors.New(config.Service.SrvName, detail, code)
}

// NewError generates a custom error.
func NewError(err error, code int32) error {
	return errors.New(config.Service.SrvName, err.Error(), code)
}

// BadRequest generates a 400 error.
func BadRequest(id, format string, a ...interface{}) error {
	return errors.BadRequest(config.Service.SrvName, format, a...)
}

// Unauthorized generates a 401 error.
func Unauthorized(format string, a ...interface{}) error {
	return errors.Unauthorized(config.Service.SrvName, format, a...)
}

// Forbidden generates a 403 error.
func Forbidden(format string, a ...interface{}) error {
	return errors.Forbidden(config.Service.SrvName, format, a...)
}

// NotFound generates a 404 error.
func NotFound(format string, a ...interface{}) error {
	return errors.NotFound(config.Service.SrvName, format, a...)
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(format string, a ...interface{}) error {
	return errors.MethodNotAllowed(config.Service.SrvName, format, a...)
}

// Timeout generates a 408 error.
func Timeout(format string, a ...interface{}) error {
	return errors.Timeout(config.Service.SrvName, format, a...)
}

// Conflict generates a 409 error.
func Conflict(format string, a ...interface{}) error {
	return errors.Conflict(config.Service.SrvName, format, a...)
}

// ServerError generates a 500 error.
func ServerError(format string, a ...interface{}) error {
	return errors.InternalServerError(config.Service.SrvName, format, a...)
}
