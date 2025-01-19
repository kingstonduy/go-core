package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	DefaultSuccessStatusCode      = ErrorStatusSuccessful
	DefaultSuccessResponseCode    = ErrorCodeSuccessful
	DefaultSuccessResponseMessage = ErrorMessageSuccessful

	DefaultFailureStatusCode      = ErrorStatusFailed
	DefaultFailureResponseCode    = ErrorCodeFailed
	DefaultFailureResponseMessage = ErrorMessageFailed
)

type Error struct {
	Status  int         `json:"status"`  // http status mapping
	Code    string      `json:"code"`    // error code
	Message string      `json:"message"` // error message
	Details interface{} `json:"details"` // error details
}

func (err *Error) Error() string {
	errMsg, _ := json.Marshal(err)
	return string(errMsg)
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Message = err
	}
	return e
}

func NewError(status int, code string, format string, a ...interface{}) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: strings.TrimSpace(fmt.Sprintf(format, a...)),
	}
}

func NewErrorWithDetails(status int, code string, details interface{}, format string, a ...interface{}) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: strings.TrimSpace(fmt.Sprintf(format, a...)),
		Details: details,
	}
}

// Deprecated: Used Failed with instead
func UnknownError(format string, a ...interface{}) *Error {
	return NewError(
		DefaultFailureStatusCode,
		DefaultFailureResponseCode,
		fmt.Sprintf("%s: %s", DefaultFailureResponseMessage, format),
		a...)
}

// ============== STANDARDS ERRORS =============================

const (
	ErrorCodeSuccessful          string = "00"
	ErrorCodeFailed              string = "01"
	ErrorCodeValidation          string = "02"
	ErrorCodeNotFound            string = "03"
	ErrorCodeOutbound            string = "04"
	ErrorCodeTimeout             string = "05"
	ErrorCodeBadRequest          string = "06"
	ErrorCodeUnauthorized        string = "07"
	ErrorCodeForbidden           string = "08"
	ErrorCodeMethodNotAllowed    string = "09"
	ErrorCodeConflict            string = "10"
	ErrorCodeTooManyRequests     string = "11"
	ErrorCodeNoRowAffected       string = "12"
	ErrorAuthenticationError     string = "13"
	ErrorCodeSuspendedError      string = "14"
	ErrorCodeInternalServerError string = "999"

	ErrorStatusSuccessful          int = http.StatusOK
	ErrorStatusFailed              int = http.StatusOK
	ErrorStatusValidation          int = http.StatusOK
	ErrorStatusNotFound            int = http.StatusOK
	ErrorStatusOutbound            int = http.StatusOK
	ErrorStatusTimeout             int = http.StatusOK
	ErrorStatusBadRequest          int = http.StatusOK
	ErrorStatusUnauthorized        int = http.StatusOK
	ErrorStatusForbidden           int = http.StatusOK
	ErrorStatusMethodNotAllowed    int = http.StatusOK
	ErrorStatusConflict            int = http.StatusOK
	ErrorStatusTooManyRequests     int = http.StatusOK
	ErrorStatusNoRowAffected       int = http.StatusOK
	ErrorStatusAuthentication      int = http.StatusOK
	ErrorStatusSuspendedError      int = http.StatusOK
	ErrorStatusInternalServerError int = http.StatusInternalServerError

	ErrorMessageSuccessful          string = "Successful"
	ErrorMessageFailed              string = "Failed"
	ErrorMessageValidation          string = "Validation error"
	ErrorMessageNotFound            string = "Not found error"
	ErrorMessageOutbound            string = "Outbound error"
	ErrorMessageTimeout             string = "Timeout error"
	ErrorMessageBadRequest          string = "Bad request error"
	ErrorMessageUnauthorized        string = "Unauthorized error"
	ErrorMessageForbidden           string = "Forbidden error"
	ErrorMessageMethodNotAllowed    string = "Method not allowed error"
	ErrorMessageConflict            string = "Conflict error"
	ErrorMessageTooManyRequests     string = "Too many request error"
	ErrorMessageNoRowAffected       string = "No row affected"
	ErrorMessagesAuthentication     string = "authentication error"
	ErrorMessageSuspendedError      string = "suspended error"
	ErrorMessageInternalServerError string = "Internal server error"
)

// Failed error
func Failed(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageFailed, format, a...)
	return NewError(ErrorStatusFailed, ErrorCodeFailed, message)
}

// Failed error with details
func FailedWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageFailed, format, a...)
	return NewErrorWithDetails(ErrorStatusFailed, ErrorCodeFailed, details, message)
}

// Validation error
func ValidationError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageValidation, format, a...)
	return NewError(ErrorStatusValidation, ErrorCodeValidation, message)
}

// Validation error with details
func ValidationErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageValidation, format, a...)
	return NewErrorWithDetails(ErrorStatusValidation, ErrorCodeValidation, details, message)
}

// Not found error
func NotFoundError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageNotFound, format, a...)
	return NewError(ErrorStatusNotFound, ErrorCodeNotFound, message)
}

// Not found error with details
func NotFoundErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageNotFound, format, a...)
	return NewErrorWithDetails(ErrorStatusNotFound, ErrorCodeNotFound, details, message)
}

// Outbound error
func OutboundError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageOutbound, format, a...)
	return NewError(ErrorStatusOutbound, ErrorCodeOutbound, message)
}

// Outbound error with details
func OutboundErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageOutbound, format, a...)
	return NewErrorWithDetails(ErrorStatusOutbound, ErrorCodeOutbound, details, message)
}

// Timeout error
func TimeoutError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageTimeout, format, a...)
	return NewError(ErrorStatusTimeout, ErrorCodeTimeout, message)
}

// Timeout error with details
func TimeoutErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageTimeout, format, a...)
	return NewErrorWithDetails(ErrorStatusTimeout, ErrorCodeTimeout, details, message)
}

// Bad request error
func BadRequestError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageBadRequest, format, a...)
	return NewError(ErrorStatusBadRequest, ErrorCodeBadRequest, message)
}

// Bad request error with details
func BadRequestErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageBadRequest, format, a...)
	return NewErrorWithDetails(ErrorStatusBadRequest, ErrorCodeBadRequest, details, message)
}

// Unauthorized error
func UnauthorizedError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageUnauthorized, format, a...)
	return NewError(ErrorStatusUnauthorized, ErrorCodeUnauthorized, message)
}

// Unauthorized error with details
func UnauthorizedErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageUnauthorized, format, a...)
	return NewErrorWithDetails(ErrorStatusUnauthorized, ErrorCodeUnauthorized, details, message)
}

// Forbidden error
func ForbiddenError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageForbidden, format, a...)
	return NewError(ErrorStatusForbidden, ErrorCodeForbidden, message)
}

// Forbidden error with details
func ForbiddenErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageForbidden, format, a...)
	return NewErrorWithDetails(ErrorStatusForbidden, ErrorCodeForbidden, details, message)
}

// Method not allowed error
func MethodNotAllowedError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageMethodNotAllowed, format, a...)
	return NewError(ErrorStatusMethodNotAllowed, ErrorCodeMethodNotAllowed, message)
}

// Method not allowed error with details
func MethodNotAllowedErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageMethodNotAllowed, format, a...)
	return NewErrorWithDetails(ErrorStatusMethodNotAllowed, ErrorCodeMethodNotAllowed, details, message)
}

// Conflict error
func ConflictError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageConflict, format, a...)
	return NewError(ErrorStatusConflict, ErrorCodeConflict, message)
}

// Conflict error with details
func ConflictErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageConflict, format, a...)
	return NewErrorWithDetails(ErrorStatusConflict, ErrorCodeConflict, details, message)
}

// Too many request error
func TooManyRequestError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageTooManyRequests, format, a...)
	return NewError(ErrorStatusTooManyRequests, ErrorCodeTooManyRequests, message)
}

// Too many request error with details
func TooManyRequestErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageTooManyRequests, format, a...)
	return NewErrorWithDetails(ErrorStatusTooManyRequests, ErrorCodeTooManyRequests, details, message)
}

// No row affected error
func NoRowAffectedError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageNoRowAffected, format, a...)
	return NewError(ErrorStatusNoRowAffected, ErrorCodeNoRowAffected, message)
}

// No row affected error with details
func NoRowAffectedErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageNoRowAffected, format, a...)
	return NewErrorWithDetails(ErrorStatusNoRowAffected, ErrorCodeNoRowAffected, details, message)
}

func AuthenticationError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessagesAuthentication, format, a...)
	return NewError(ErrorStatusAuthentication, ErrorAuthenticationError, message)
}

// authentication error with details
func AuthenticationErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessagesAuthentication, format, a...)
	return NewErrorWithDetails(ErrorStatusAuthentication, ErrorAuthenticationError, details, message)
}

// suspended error
func SuspendedError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageSuspendedError, format, a...)
	return NewError(ErrorStatusSuspendedError, ErrorCodeSuspendedError, message)
}

// suspended error with details
func SuspendedErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageSuspendedError, format, a...)
	return NewErrorWithDetails(ErrorStatusSuspendedError, ErrorCodeSuspendedError, details, message)
}

// Internal server error
func InternalServerError(format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageInternalServerError, format, a...)
	return NewError(ErrorStatusInternalServerError, ErrorCodeInternalServerError, message)
}

// Internal server error with details
func InternalServerErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	message := buildErrorMessage(ErrorMessageInternalServerError, format, a...)
	return NewErrorWithDetails(ErrorStatusInternalServerError, ErrorCodeInternalServerError, details, message)
}

// Helper function to build error messages
func buildErrorMessage(baseMessage, format string, a ...interface{}) string {
	if format == "" {
		return baseMessage
	}
	return fmt.Sprintf("%s. %s", baseMessage, fmt.Sprintf(format, a...))
}

func Equal(err1 error, err2 error) bool {
	verr1, ok1 := err1.(*Error)
	verr2, ok2 := err2.(*Error)

	if ok1 != ok2 {
		return false
	}

	if !ok1 {
		return err1 == err2
	}

	if verr1.Status != verr2.Status {
		return false
	}

	if verr1.Code != verr2.Code {
		return false
	}

	return true
}

// FromError try to convert go error to *Error.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if verr, ok := err.(*Error); ok && verr != nil {
		return verr
	}

	return Parse(err.Error())
}

// As finds the first error in err's chain that matches *Error.
func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var merr *Error
	if errors.As(err, &merr) {
		return merr, true
	}
	return nil, false
}
