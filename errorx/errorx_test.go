package errorx

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorx(t *testing.T) {
	t.Run("TestError", TestError)
	t.Run("TestErrorDetail", TestErrorDetail)
	t.Run("TestParse", TestParse)
	t.Run("TestEqual", TestEqual)
	t.Run("TestFromError", TestFromError)
	t.Run("TestAs", TestAs)
	t.Run("TestErrorTypes", TestErrorsTypes)
}

func TestErrorsTypes(t *testing.T) {
	tests := []struct {
		name      string
		errorFunc func() *Error
		expected  struct {
			status  int
			code    string
			message string
			details interface{}
		}
	}{
		{
			name:      "TestFailedError",
			errorFunc: func() *Error { return Failed("Something failed") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusFailed,
				code:    ErrorCodeFailed,
				message: ErrorMessageFailed + ". Something failed",
				details: nil,
			},
		},
		{
			name:      "TestValidationError",
			errorFunc: func() *Error { return ValidationError("Invalid input") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusValidation,
				code:    ErrorCodeValidation,
				message: ErrorMessageValidation + ". Invalid input",
				details: nil,
			},
		},
		{
			name:      "TestNotFoundError",
			errorFunc: func() *Error { return NotFoundError("Resource not found") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusNotFound,
				code:    ErrorCodeNotFound,
				message: ErrorMessageNotFound + ". Resource not found",
				details: nil,
			},
		},
		{
			name:      "TestOutboundError",
			errorFunc: func() *Error { return OutboundError("Failed to reach external service") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusOutbound,
				code:    ErrorCodeOutbound,
				message: ErrorMessageOutbound + ". Failed to reach external service",
				details: nil,
			},
		},
		{
			name:      "TestTimeoutError",
			errorFunc: func() *Error { return TimeoutError("Request timed out") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusTimeout,
				code:    ErrorCodeTimeout,
				message: ErrorMessageTimeout + ". Request timed out",
				details: nil,
			},
		},
		{
			name:      "TestBadRequestError",
			errorFunc: func() *Error { return BadRequestError("Bad request") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusBadRequest,
				code:    ErrorCodeBadRequest,
				message: ErrorMessageBadRequest + ". Bad request",
				details: nil,
			},
		},
		{
			name:      "TestUnauthorizedError",
			errorFunc: func() *Error { return UnauthorizedError("Unauthorized access") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusUnauthorized,
				code:    ErrorCodeUnauthorized,
				message: ErrorMessageUnauthorized + ". Unauthorized access",
				details: nil,
			},
		},
		{
			name:      "TestForbiddenError",
			errorFunc: func() *Error { return ForbiddenError("Access denied") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusForbidden,
				code:    ErrorCodeForbidden,
				message: ErrorMessageForbidden + ". Access denied",
				details: nil,
			},
		},
		{
			name:      "TestMethodNotAllowedError",
			errorFunc: func() *Error { return MethodNotAllowedError("Method not allowed") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusMethodNotAllowed,
				code:    ErrorCodeMethodNotAllowed,
				message: ErrorMessageMethodNotAllowed + ". Method not allowed",
				details: nil,
			},
		},
		{
			name:      "TestConflictError",
			errorFunc: func() *Error { return ConflictError("Conflict occurred") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusConflict,
				code:    ErrorCodeConflict,
				message: ErrorMessageConflict + ". Conflict occurred",
				details: nil,
			},
		},
		{
			name:      "TestTooManyRequestError",
			errorFunc: func() *Error { return TooManyRequestError("Too many requests") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusTooManyRequests,
				code:    ErrorCodeTooManyRequests,
				message: ErrorMessageTooManyRequests + ". Too many requests",
				details: nil,
			},
		},
		{
			name:      "TestInternalServerError",
			errorFunc: func() *Error { return InternalServerError("Internal server error occurred") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusInternalServerError,
				code:    ErrorCodeInternalServerError,
				message: ErrorMessageInternalServerError + ". Internal server error occurred",
				details: nil,
			},
		},
		{
			name: "TestFailedWithDetails",
			errorFunc: func() *Error {
				return FailedWithDetails(map[string]interface{}{"context": "test"}, "Failed due to an issue")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusFailed,
				code:    ErrorCodeFailed,
				message: ErrorMessageFailed + ". Failed due to an issue",
				details: map[string]interface{}{"context": "test"},
			},
		},
		{
			name: "TestValidationErrorWithDetails",
			errorFunc: func() *Error {
				return ValidationErrorWithDetails(map[string]interface{}{"field": "email", "error": "invalid format"}, "Email validation error")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusValidation,
				code:    ErrorCodeValidation,
				message: ErrorMessageValidation + ". Email validation error",
				details: map[string]interface{}{"field": "email", "error": "invalid format"},
			},
		},
		{
			name: "TestNotFoundErrorWithDetails",
			errorFunc: func() *Error {
				return NotFoundErrorWithDetails(map[string]interface{}{"resource": "user"}, "User not found")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusNotFound,
				code:    ErrorCodeNotFound,
				message: ErrorMessageNotFound + ". User not found",
				details: map[string]interface{}{"resource": "user"},
			},
		},
		{
			name: "TestOutboundErrorWithDetails",
			errorFunc: func() *Error {
				return OutboundErrorWithDetails(map[string]interface{}{"service": "payment gateway"}, "Failed to connect")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusOutbound,
				code:    ErrorCodeOutbound,
				message: ErrorMessageOutbound + ". Failed to connect",
				details: map[string]interface{}{"service": "payment gateway"},
			},
		},
		{
			name: "TestTimeoutErrorWithDetails",
			errorFunc: func() *Error {
				return TimeoutErrorWithDetails(map[string]interface{}{"timeout": "30s"}, "Request timed out")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusTimeout,
				code:    ErrorCodeTimeout,
				message: ErrorMessageTimeout + ". Request timed out",
				details: map[string]interface{}{"timeout": "30s"},
			},
		},
		{
			name: "TestBadRequestErrorWithDetails",
			errorFunc: func() *Error {
				return BadRequestErrorWithDetails(map[string]interface{}{"parameter": "username"}, "Bad request")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusBadRequest,
				code:    ErrorCodeBadRequest,
				message: ErrorMessageBadRequest + ". Bad request",
				details: map[string]interface{}{"parameter": "username"},
			},
		},
		{
			name: "TestUnauthorizedErrorWithDetails",
			errorFunc: func() *Error {
				return UnauthorizedErrorWithDetails(map[string]interface{}{"token": "expired"}, "Unauthorized access")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusUnauthorized,
				code:    ErrorCodeUnauthorized,
				message: ErrorMessageUnauthorized + ". Unauthorized access",
				details: map[string]interface{}{"token": "expired"},
			},
		},
		{
			name: "TestForbiddenErrorWithDetails",
			errorFunc: func() *Error {
				return ForbiddenErrorWithDetails(map[string]interface{}{"user": "guest"}, "Access denied")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusForbidden,
				code:    ErrorCodeForbidden,
				message: ErrorMessageForbidden + ". Access denied",
				details: map[string]interface{}{"user": "guest"},
			},
		},
		{
			name: "TestMethodNotAllowedErrorWithDetails",
			errorFunc: func() *Error {
				return MethodNotAllowedErrorWithDetails(map[string]interface{}{"method": "POST"}, "Method not allowed")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusMethodNotAllowed,
				code:    ErrorCodeMethodNotAllowed,
				message: ErrorMessageMethodNotAllowed + ". Method not allowed",
				details: map[string]interface{}{"method": "POST"},
			},
		},
		{
			name: "TestConflictErrorWithDetails",
			errorFunc: func() *Error {
				return ConflictErrorWithDetails(map[string]interface{}{"resource": "email"}, "Conflict occurred")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusConflict,
				code:    ErrorCodeConflict,
				message: ErrorMessageConflict + ". Conflict occurred",
				details: map[string]interface{}{"resource": "email"},
			},
		},
		{
			name: "TestTooManyRequestErrorWithDetails",
			errorFunc: func() *Error {
				return TooManyRequestErrorWithDetails(map[string]interface{}{"retryAfter": "60s"}, "Too many requests")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusTooManyRequests,
				code:    ErrorCodeTooManyRequests,
				message: ErrorMessageTooManyRequests + ". Too many requests",
				details: map[string]interface{}{"retryAfter": "60s"},
			},
		},
		{
			name: "TestInternalServerErrorWithDetails",
			errorFunc: func() *Error {
				return InternalServerErrorWithDetails(map[string]interface{}{"service": "database"}, "Internal server error occurred")
			},
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusInternalServerError,
				code:    ErrorCodeInternalServerError,
				message: ErrorMessageInternalServerError + ". Internal server error occurred",
				details: map[string]interface{}{"service": "database"},
			},
		},
		{
			name:      "TestFailedErrorWithEmptyMessage",
			errorFunc: func() *Error { return Failed("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusFailed,
				code:    ErrorCodeFailed,
				message: ErrorMessageFailed,
				details: nil,
			},
		},
		{
			name:      "TestValidationErrorWithEmptyMessage",
			errorFunc: func() *Error { return ValidationError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusValidation,
				code:    ErrorCodeValidation,
				message: ErrorMessageValidation,
				details: nil,
			},
		},
		{
			name:      "TestNotFoundErrorWithEmptyMessage",
			errorFunc: func() *Error { return NotFoundError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusNotFound,
				code:    ErrorCodeNotFound,
				message: ErrorMessageNotFound,
				details: nil,
			},
		},
		{
			name:      "TestOutboundErrorWithEmptyMessage",
			errorFunc: func() *Error { return OutboundError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusOutbound,
				code:    ErrorCodeOutbound,
				message: ErrorMessageOutbound,
				details: nil,
			},
		},
		{
			name:      "TestTimeoutErrorWithEmptyMessage",
			errorFunc: func() *Error { return TimeoutError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusTimeout,
				code:    ErrorCodeTimeout,
				message: ErrorMessageTimeout,
				details: nil,
			},
		},
		{
			name:      "TestBadRequestErrorWithEmptyMessage",
			errorFunc: func() *Error { return BadRequestError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusBadRequest,
				code:    ErrorCodeBadRequest,
				message: ErrorMessageBadRequest,
				details: nil,
			},
		},
		{
			name:      "TestUnauthorizedErrorWithEmptyMessage",
			errorFunc: func() *Error { return UnauthorizedError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusUnauthorized,
				code:    ErrorCodeUnauthorized,
				message: ErrorMessageUnauthorized,
				details: nil,
			},
		},
		{
			name:      "TestForbiddenErrorWithEmptyMessage",
			errorFunc: func() *Error { return ForbiddenError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusForbidden,
				code:    ErrorCodeForbidden,
				message: ErrorMessageForbidden,
				details: nil,
			},
		},
		{
			name:      "TestMethodNotAllowedErrorWithEmptyMessage",
			errorFunc: func() *Error { return MethodNotAllowedError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusMethodNotAllowed,
				code:    ErrorCodeMethodNotAllowed,
				message: ErrorMessageMethodNotAllowed,
				details: nil,
			},
		},
		{
			name:      "TestConflictErrorWithEmptyMessage",
			errorFunc: func() *Error { return ConflictError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusConflict,
				code:    ErrorCodeConflict,
				message: ErrorMessageConflict,
				details: nil,
			},
		},
		{
			name:      "TestTooManyRequestErrorWithEmptyMessage",
			errorFunc: func() *Error { return TooManyRequestError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusTooManyRequests,
				code:    ErrorCodeTooManyRequests,
				message: ErrorMessageTooManyRequests,
				details: nil,
			},
		},
		{
			name:      "TestInternalServerErrorWithEmptyMessage",
			errorFunc: func() *Error { return InternalServerError("") },
			expected: struct {
				status  int
				code    string
				message string
				details interface{}
			}{
				status:  ErrorStatusInternalServerError,
				code:    ErrorCodeInternalServerError,
				message: ErrorMessageInternalServerError,
				details: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errorFunc()
			assert.Equal(t, tt.expected.status, err.Status)
			assert.Equal(t, tt.expected.code, err.Code)
			assert.Equal(t, tt.expected.message, err.Message)
			assert.Equal(t, tt.expected.details, err.Details)
		})
	}
}

func TestError(t *testing.T) {
	err := NewError(99, "99", "test: %s", "data")

	assert.Nil(t, err.Details)
	assert.Equal(t, 99, err.Status)
	assert.Equal(t, "99", err.Code)
	assert.Equal(t, "test: data", err.Message)
}

func TestErrorDetail(t *testing.T) {
	details := []string{"foo", "bar"}
	err := NewErrorWithDetails(99, "99", details, "test: %s", "data")

	assert.Equal(t, details, err.Details)
	assert.Equal(t, 99, err.Status)
	assert.Equal(t, "99", err.Code)
	assert.Equal(t, "test: data", err.Message)
}

func TestParse(t *testing.T) {
	jsonErr := `{"status": 404, "code": "03", "message": "Not found", "details": null}`
	err := Parse(jsonErr)

	assert.Equal(t, 404, err.Status)
	assert.Equal(t, "03", err.Code)
	assert.Equal(t, "Not found", err.Message)
	assert.Nil(t, err.Details)

	// Test parsing invalid JSON
	invalidJsonErr := `{"status": 404, "code": "03", "message": "Not found"`
	err = Parse(invalidJsonErr)

	assert.Equal(t, invalidJsonErr, err.Message) // should fall back to original string
	assert.Equal(t, 0, err.Status)
	assert.Empty(t, err.Code)
}

func TestEqual(t *testing.T) {
	err1 := NewError(http.StatusBadRequest, "400", "This is a test error")
	err2 := NewError(http.StatusBadRequest, "400", "This is a test error")
	err3 := NewError(http.StatusInternalServerError, "999", "Internal server error")

	assert.True(t, Equal(err1, err2))
	assert.False(t, Equal(err1, err3))
	assert.False(t, Equal(err1, nil))
}

func TestFromError(t *testing.T) {
	err := NewError(http.StatusBadRequest, "400", "This is a test error")
	result := FromError(err)

	assert.NotNil(t, result)
	assert.Equal(t, err, result)

	// Test conversion from a standard error
	stdErr := errors.New("standard error")
	result = FromError(stdErr)

	assert.NotNil(t, result)
	assert.Equal(t, stdErr.Error(), result.Message)

	// Test nil error
	result = FromError(nil)
	assert.Nil(t, result)
}

func TestAs(t *testing.T) {
	err := NewError(http.StatusNotFound, "03", "Not found error")
	result, ok := As(err)

	assert.True(t, ok)
	assert.Equal(t, err, result)

	// Test non-*Error
	stdErr := errors.New("standard error")
	result, ok = As(stdErr)

	assert.False(t, ok)
	assert.Nil(t, result)
}
