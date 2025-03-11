package errors

import (
	"encoding/json"
)

// ErrorType is the list of allowed values for the error's type.
type ErrorType string

// List of values that ErrorType can take.
const (
	ErrorTypeAPI            ErrorType = "api_error"
	ErrorTypeInvalidRequest ErrorType = "invalid_request_error"
)

// ErrorCode is the list of allowed values for the error's code.
type ErrorCode string

const (
	ErrorCodeRateLimit                   ErrorCode = "rate_limit"
	ErrorCodeNotAllowedOnStandardAccount ErrorCode = "not_allowed_on_standard_account"
	ErrorCodeParameterInvalidDate        ErrorCode = "parameter_invalid_date"
	ErrorCodeParameterInvalidEmpty       ErrorCode = "parameter_invalid_empty"
	ErrorCodeParameterInvalidInteger     ErrorCode = "parameter_invalid_integer"
	ErrorCodeParameterInvalidStringEmpty ErrorCode = "parameter_invalid_string_empty"
	ErrorCodeParameterMissing            ErrorCode = "parameter_missing"
	ErrorCodeNotFound                    ErrorCode = "not_found"
	ErrorCodeResourceAlreadyExists       ErrorCode = "resource_already_exists"
	ErrorCodeResourceMissing             ErrorCode = "resource_missing"
)

// Error is the response returned when a call is unsuccessful.
type Error struct {
	// The HTTP status code of the error.
	Status int `json:"status"`
	// A machine-readable code indicating the type of error.
	Code string `json:"code"`
	// A summary of the error.
	Title string `json:"title,omitempty"`
	// A detailed explanation of the error.
	Detail string `json:"detail,omitempty"`
	Source *struct {
		// The query parameter that produced the error.
		Parameter string `json:"parameter,omitempty"`
		// A JSON pointer that indicates the location in the request entity where the error originates.
		Pointer string `json:"pointer,omitempty"`
	} `json:"source,omitempty"`
}

// Error serializes the error object to JSON and returns it as a string.
func (e Error) Error() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}

// APIError is a catch-all for any errors not covered by other types
type APIError struct {
	error *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *APIError) Error() string {
	return e.error.Error()
}

// InvalidRequestError is an error that occurs when a request contains invalid parameters.
type InvalidRequestError struct {
	error *Error
}

// Error serializes the error object to JSON and returns it as a string.
func (e *InvalidRequestError) Error() string {
	return e.error.Error()
}

type ErrorsResponse struct {
	Errors []Error `json:"errors"`
}

type ForbiddenResponse struct {
	Errors []Error
}

func NewHTTPError(status int, code ErrorCode, title string, details ...string) *Error {
	error := &Error{
		Status: status,
		Code:   string(code),
		Title:  title,
	}
	if len(details) > 0 {
		error.Detail = details[0]
	}
	return error
}
