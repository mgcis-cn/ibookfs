// Package errors provides a unified error handling framework.
//
// It defines a standard Error interface that carries an application error code,
// HTTP status code, and human-readable message. The pkg layer provides the
// interface and base implementation; application modules (e.g. apiserver) define
// their own error codes and sentinel errors on top of this.
//
// Usage:
//
//	// Define module errors (in internal/apiserver/errors/)
//	var ErrUserNotFound = errors.New(10001, http.StatusNotFound, "用户不存在")
//
//	// Return in business logic
//	return nil, ErrUserNotFound
//
//	// Wrap with extra context
//	return nil, errors.WithMessage(ErrUserNotFound, "email not registered")
//
//	// Extract in error encoder
//	if e := errors.FromError(err); e != nil {
//	    // use e.Code(), e.HTTPStatus(), e.Message()
//	}
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Error is the unified error interface for the application.
// It extends the standard error interface with structured fields.
type Error interface {
	error
	// Code returns the application-level error code.
	Code() int
	// HTTPStatus returns the HTTP status code to respond with.
	HTTPStatus() int
	// Message returns the user-facing error message.
	Message() string
	// Unwrap returns the underlying error, if any.
	Unwrap() error
}

// appError is the concrete implementation of Error.
type appError struct {
	code       int
	httpStatus int
	message    string
	cause      error
}

// New creates a new Error with the given code, HTTP status, and message.
func New(code int, httpStatus int, message string) Error {
	return &appError{
		code:       code,
		httpStatus: httpStatus,
		message:    message,
	}
}

// Newf creates a new Error with a formatted message.
func Newf(code int, httpStatus int, format string, args ...interface{}) Error {
	return &appError{
		code:       code,
		httpStatus: httpStatus,
		message:    fmt.Sprintf(format, args...),
	}
}

// WithMessage wraps an existing Error with a new message, preserving code and HTTP status.
// If err is not an Error, it wraps as an internal server error.
func WithMessage(err error, message string) Error {
	if e := FromError(err); e != nil {
		return &appError{
			code:       e.Code(),
			httpStatus: e.HTTPStatus(),
			message:    message,
			cause:      err,
		}
	}
	return &appError{
		code:       0,
		httpStatus: http.StatusInternalServerError,
		message:    message,
		cause:      err,
	}
}

// WithCause wraps an existing Error with an underlying cause error.
// The original Error's code, HTTP status, and message are preserved.
func WithCause(err Error, cause error) Error {
	return &appError{
		code:       err.Code(),
		httpStatus: err.HTTPStatus(),
		message:    err.Message(),
		cause:      cause,
	}
}

// FromError extracts an Error from the error chain.
// Returns nil if err is nil or does not contain an Error.
func FromError(err error) Error {
	if err == nil {
		return nil
	}
	var e *appError
	if errors.As(err, &e) {
		return e
	}
	return nil
}

// IsCode checks whether the error chain contains an Error with the given code.
func IsCode(err error, code int) bool {
	if e := FromError(err); e != nil {
		return e.Code() == code
	}
	return false
}

// --- Error interface implementation ---

func (e *appError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *appError) Code() int       { return e.code }
func (e *appError) HTTPStatus() int  { return e.httpStatus }
func (e *appError) Message() string  { return e.message }
func (e *appError) Unwrap() error    { return e.cause }

// --- Common HTTP error constructors ---

// BadRequest creates a 400 Bad Request error.
func BadRequest(code int, message string) Error {
	return New(code, http.StatusBadRequest, message)
}

// Unauthorized creates a 401 Unauthorized error.
func Unauthorized(code int, message string) Error {
	return New(code, http.StatusUnauthorized, message)
}

// Forbidden creates a 403 Forbidden error.
func Forbidden(code int, message string) Error {
	return New(code, http.StatusForbidden, message)
}

// NotFound creates a 404 Not Found error.
func NotFound(code int, message string) Error {
	return New(code, http.StatusNotFound, message)
}

// Conflict creates a 409 Conflict error.
func Conflict(code int, message string) Error {
	return New(code, http.StatusConflict, message)
}

// TooManyRequests creates a 429 Too Many Requests error.
func TooManyRequests(code int, message string) Error {
	return New(code, http.StatusTooManyRequests, message)
}

// InternalServer creates a 500 Internal Server Error.
func InternalServer(code int, message string) Error {
	return New(code, http.StatusInternalServerError, message)
}

// --- Re-export standard library functions for convenience ---

// Is reports whether any error in err's tree matches target.
func Is(err, target error) bool { return errors.Is(err, target) }

// As finds the first error in err's tree that matches target.
func As(err error, target interface{}) bool { return errors.As(err, target) }

// Wrap wraps a standard error as an internal server error with a message.
func Wrap(err error, message string) Error {
	return &appError{
		code:       0,
		httpStatus: http.StatusInternalServerError,
		message:    message,
		cause:      err,
	}
}
