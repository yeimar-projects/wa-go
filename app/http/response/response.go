package response

import (
	"log/slog"

	contractshttp "github.com/goravel/framework/contracts/http"

	apperrors "github.com/yeimar-sandbox/wa-go/app/errors"
)

// ---------------------------------------------------------------------------
// Standard JSON envelope — every API response uses this shape.
// ---------------------------------------------------------------------------

// Link represents a HATEOAS link
type Link struct {
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}

// Links is a map of relation name to Link
type Links map[string]Link

// Cursor pagination
type Pagination struct {
	Limit   int     `json:"limit"`
	HasMore bool    `json:"hasMore"`
	Cursors Cursors `json:"cursors"`
}

type Cursors struct {
	After  string `json:"after,omitempty"`
	Before string `json:"before,omitempty"`
}

// ---------------------------------------------------------------------------
// Unified response structs
// ---------------------------------------------------------------------------

// APIResponse is the single canonical JSON envelope for ALL API responses.
type APIResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Results any    `json:"results,omitempty"`
}

// PaginatedResponse extends APIResponse with pagination metadata.
type PaginatedResponse struct {
	Status     int        `json:"status"`
	Code       string     `json:"code"`
	Message    string     `json:"message"`
	Results    any        `json:"results"`
	Pagination Pagination `json:"pagination"`
}

// ErrorDetail provides machine-readable error information.
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// ErrorResponse is the canonical error envelope.
type ErrorResponse struct {
	Status  int           `json:"status"`
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

// ---------------------------------------------------------------------------
// Success helpers
// ---------------------------------------------------------------------------

// NewSuccess creates a standard success response with HTTP 200.
func NewSuccess(results any, message string) APIResponse {
	return APIResponse{
		Status:  200,
		Code:    "SUCCESS",
		Message: message,
		Results: results,
	}
}

// NewCreated creates a success response with HTTP 201.
func NewCreated(results any, message string) APIResponse {
	return APIResponse{
		Status:  201,
		Code:    "CREATED",
		Message: message,
		Results: results,
	}
}

// NewPaginated creates a paginated success response.
func NewPaginated(results any, limit int, hasMore bool, after, before string, message string) PaginatedResponse {
	return PaginatedResponse{
		Status:  200,
		Code:    "SUCCESS",
		Message: message,
		Results: results,
		Pagination: Pagination{
			Limit:   limit,
			HasMore: hasMore,
			Cursors: Cursors{After: after, Before: before},
		},
	}
}

// ---------------------------------------------------------------------------
// Error helpers
// ---------------------------------------------------------------------------

// NewError creates an error response from a status code, machine code, and message.
// This is the low-level constructor; prefer Error() for AppError integration.
func NewError(status int, code, message string) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

// NewValidationError creates a 400 error with field-level details.
func NewValidationError(message string, details []ErrorDetail) ErrorResponse {
	return ErrorResponse{
		Status:  400,
		Code:    string(apperrors.CodeValidation),
		Message: message,
		Errors:  details,
	}
}

// ---------------------------------------------------------------------------
// AppError integration — the recommended way to handle errors in controllers.
// ---------------------------------------------------------------------------

// Error converts any error into a proper HTTP JSON response.
// If the error is an *apperrors.AppError, its code and message are used.
// Otherwise, a generic 500 is returned (internal details are logged, not leaked).
//
// Usage in controllers:
//
//	if err != nil {
//	    return response.Error(ctx, err)
//	}
func Error(ctx contractshttp.Context, err error) contractshttp.Response {
	appErr := apperrors.ToAppError(err)

	// Log the full error chain for observability (structured logging).
	slog.Error("request failed",
		"code", string(appErr.Code),
		"message", appErr.Message,
		"internal", appErr.Internal,
		"path", ctx.Request().Path(),
		"method", ctx.Request().Method(),
	)

	status := appErr.HTTPStatus()
	return ctx.Response().Json(status, ErrorResponse{
		Status:  status,
		Code:    string(appErr.Code),
		Message: appErr.Message,
	})
}

// ---------------------------------------------------------------------------
// HATEOAS link helpers
// ---------------------------------------------------------------------------

// Helper to build instance links based on status
func InstanceLinks(baseURL, id, status string) Links {
	self := baseURL + "/instances/" + id
	links := Links{
		"self": {Href: self, Method: "GET"},
	}
	switch status {
	case "connected":
		links["disconnect"] = Link{Href: self + "/disconnect", Method: "POST"}
		links["logout"] = Link{Href: self + "/logout", Method: "POST"}
		links["messages"] = Link{Href: self + "/messages", Method: "POST"}
		links["chats"] = Link{Href: self + "/chats", Method: "GET"}
		links["groups"] = Link{Href: self + "/groups", Method: "GET"}
		links["contacts"] = Link{Href: self + "/contacts", Method: "GET"}
		links["newsletters"] = Link{Href: self + "/newsletters", Method: "GET"}
		links["presence"] = Link{Href: self + "/presence", Method: "PUT"}
		links["privacy"] = Link{Href: self + "/privacy", Method: "GET"}
		links["profile"] = Link{Href: self + "/profile", Method: "GET"}
		links["webhooks"] = Link{Href: self + "/webhooks", Method: "GET"}
	case "disconnected":
		links["connect"] = Link{Href: self + "/connect", Method: "POST"}
		links["pair-phone"] = Link{Href: self + "/pair-phone", Method: "POST"}
		links["delete"] = Link{Href: self, Method: "DELETE"}
	case "qr_code":
		links["qr-code"] = Link{Href: self + "/qr-code", Method: "GET"}
		links["disconnect"] = Link{Href: self + "/disconnect", Method: "POST"}
	}
	return links
}
