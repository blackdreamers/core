package errors

import (
	errs "errors"

	"go-micro.dev/v4/errors"

	"github.com/blackdreamers/core/config"
)

// New generates a custom error.
func New(detail string, code int32) error {
	return errors.New(config.Service.SrvName, detail, code)
}

// NewError generates a custom error.
func NewError(err error, code int32) error {
	return errors.New(config.Service.SrvName, err.Error(), code)
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *errors.Error {
	return errors.Parse(err)
}

// BadRequest generates a 400 error.
func BadRequest(format string, a ...interface{}) error {
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

// Equal tries to compare errors
func Equal(err1, err2 error) bool {
	return errors.Equal(err1, err2)
}

// FromError try to convert go error to *Error
func FromError(err error) *errors.Error {
	return errors.FromError(err)
}

func V(err1, err2 error) bool {
	return errs.Is(err1, err2)
}
