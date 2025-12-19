package core

import "errors"

// Error represents an HTTP error
type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewError creates a new HTTP error
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// WrapError wraps an existing error with an HTTP error
func WrapError(code int, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common HTTP errors
var (
	ErrNotFound            = NewError(404, "Not Found")
	ErrBadRequest          = NewError(400, "Bad Request")
	ErrUnauthorized        = NewError(401, "Unauthorized")
	ErrForbidden           = NewError(403, "Forbidden")
	ErrInternalServerError = NewError(500, "Internal Server Error")
)

// IsHTTPError checks if an error is an HTTP error
func IsHTTPError(err error) bool {
	var httpErr *Error
	return errors.As(err, &httpErr)
}

// GetHTTPError extracts HTTP error from an error
func GetHTTPError(err error) *Error {
	var httpErr *Error
	if errors.As(err, &httpErr) {
		return httpErr
	}
	return nil
}

