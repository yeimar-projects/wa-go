// Package errors provides a structured, domain-driven error system.
//
// Design principles (inspired by Google Cloud APIs, Stripe, Twilio):
//   - Every error has a machine-readable Code (enum) and a human-readable Message.
//   - Errors carry the correct HTTP status so controllers never hard-code it.
//   - Internal details are captured for logging but NEVER leaked to the client.
//   - Sentinel errors allow services to return typed errors without importing HTTP.
//   - Wrapping preserves the full causal chain for observability.
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ---------------------------------------------------------------------------
// Error codes — machine-readable identifiers returned in every API response.
// ---------------------------------------------------------------------------

type Code string

const (
	// Client errors (4xx)
	CodeValidation      Code = "VALIDATION_ERROR"
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeForbidden       Code = "FORBIDDEN"
	CodeNotFound        Code = "NOT_FOUND"
	CodeConflict        Code = "CONFLICT"
	CodeTooManyRequests Code = "RATE_LIMIT_EXCEEDED"
	CodeNotImplemented  Code = "NOT_IMPLEMENTED"
	CodeUnprocessable   Code = "UNPROCESSABLE_ENTITY"

	// Server / infrastructure errors (5xx)
	CodeInternal    Code = "INTERNAL_ERROR"
	CodeServiceDown Code = "SERVICE_UNAVAILABLE"
	CodeTimeout     Code = "TIMEOUT"

	// Domain-specific errors
	CodeNotConnected     Code = "WHATSAPP_NOT_CONNECTED"
	CodeSendFailed       Code = "SEND_FAILED"
	CodeConnectionFailed Code = "CONNECTION_FAILED"
	CodeInvalidJID       Code = "INVALID_JID"
	CodeMediaFetchFailed Code = "MEDIA_FETCH_FAILED"
	CodeUploadFailed     Code = "UPLOAD_FAILED"
)

// codeHTTPStatus maps every Code to its canonical HTTP status.
var codeHTTPStatus = map[Code]int{
	CodeValidation:       http.StatusBadRequest,
	CodeBadRequest:       http.StatusBadRequest,
	CodeUnauthorized:     http.StatusUnauthorized,
	CodeForbidden:        http.StatusForbidden,
	CodeNotFound:         http.StatusNotFound,
	CodeConflict:         http.StatusConflict,
	CodeTooManyRequests:  http.StatusTooManyRequests,
	CodeNotImplemented:   http.StatusNotImplemented,
	CodeUnprocessable:    http.StatusUnprocessableEntity,
	CodeInternal:         http.StatusInternalServerError,
	CodeServiceDown:      http.StatusServiceUnavailable,
	CodeTimeout:          http.StatusGatewayTimeout,
	CodeNotConnected:     http.StatusServiceUnavailable,
	CodeSendFailed:       http.StatusBadGateway,
	CodeConnectionFailed: http.StatusBadGateway,
	CodeInvalidJID:       http.StatusBadRequest,
	CodeMediaFetchFailed: http.StatusBadGateway,
	CodeUploadFailed:     http.StatusBadGateway,
}

// HTTPStatus returns the HTTP status code for the given Code.
func (c Code) HTTPStatus() int {
	if s, ok := codeHTTPStatus[c]; ok {
		return s
	}
	return http.StatusInternalServerError
}

// ---------------------------------------------------------------------------
// AppError — the single error type that flows through the entire application.
// ---------------------------------------------------------------------------

// AppError is a structured application error.
// It satisfies the error interface and supports errors.Is / errors.As.
type AppError struct {
	// Code is the machine-readable error identifier returned to the client.
	Code Code `json:"code"`

	// Message is the human-readable message safe for API consumers.
	Message string `json:"message"`

	// Internal is the original error with full context (NEVER serialized).
	// Used exclusively for logging and debugging.
	Internal error `json:"-"`
}

// Error satisfies the error interface with a developer-friendly representation.
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap allows errors.Is and errors.As to traverse the chain.
func (e *AppError) Unwrap() error {
	return e.Internal
}

// HTTPStatus returns the appropriate HTTP status for this error.
func (e *AppError) HTTPStatus() int {
	return e.Code.HTTPStatus()
}

// ---------------------------------------------------------------------------
// Constructors — one per error category for maximum clarity at call sites.
// ---------------------------------------------------------------------------

// New creates a new AppError with the given code and user-facing message.
func New(code Code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Wrap creates a new AppError wrapping an internal cause.
// The message is what the API consumer sees; the cause is logged internally.
func Wrap(code Code, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, Internal: cause}
}

// --- Convenience constructors for common cases ---

func Validation(message string) *AppError {
	return &AppError{Code: CodeValidation, Message: message}
}

func NotFound(resource string) *AppError {
	return &AppError{Code: CodeNotFound, Message: fmt.Sprintf("%s not found", resource)}
}

func Unauthorized(message string) *AppError {
	return &AppError{Code: CodeUnauthorized, Message: message}
}

func Forbidden(message string) *AppError {
	return &AppError{Code: CodeForbidden, Message: message}
}

func Internal(message string, cause error) *AppError {
	return &AppError{Code: CodeInternal, Message: message, Internal: cause}
}

func NotConnected() *AppError {
	return &AppError{
		Code:    CodeNotConnected,
		Message: "WhatsApp instance is not connected. Please connect first.",
	}
}

func InvalidJID(raw string, cause error) *AppError {
	return &AppError{
		Code:     CodeInvalidJID,
		Message:  fmt.Sprintf("invalid JID format: %q", raw),
		Internal: cause,
	}
}

func SendFailed(cause error) *AppError {
	return &AppError{
		Code:     CodeSendFailed,
		Message:  "Failed to send message. Please try again.",
		Internal: cause,
	}
}

func ConnectionFailed(cause error) *AppError {
	return &AppError{
		Code:     CodeConnectionFailed,
		Message:  "Failed to establish WhatsApp connection.",
		Internal: cause,
	}
}

func MediaFetchFailed(cause error) *AppError {
	return &AppError{
		Code:     CodeMediaFetchFailed,
		Message:  "Failed to fetch media content.",
		Internal: cause,
	}
}

func UploadFailed(cause error) *AppError {
	return &AppError{
		Code:     CodeUploadFailed,
		Message:  "Failed to upload media to WhatsApp.",
		Internal: cause,
	}
}

func NotImplemented(feature string) *AppError {
	return &AppError{
		Code:    CodeNotImplemented,
		Message: fmt.Sprintf("%s is not yet implemented", feature),
	}
}

// ---------------------------------------------------------------------------
// Helpers for error inspection
// ---------------------------------------------------------------------------

// AsAppError extracts an *AppError from an error chain.
// Returns nil if the error is not an AppError.
func AsAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// ToAppError converts any error into an *AppError.
// If the error is already an AppError, it is returned as-is.
// Otherwise it is wrapped as an internal error (safe for logging, hidden from client).
func ToAppError(err error) *AppError {
	if appErr := AsAppError(err); appErr != nil {
		return appErr
	}
	return Internal("An unexpected error occurred. Please try again later.", err)
}
