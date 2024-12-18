package fiberx

import (
	"encoding/json"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kingstonduy/go-core/errorx"
	"github.com/stretchr/testify/assert"
)

func TestWrapFiberError(t *testing.T) {
	tests := []struct {
		name        string
		fiberError  *fiber.Error
		expectedErr error
	}{
		{
			name:        "BadRequest",
			fiberError:  fiber.NewError(fiber.StatusBadRequest, "Bad Request"),
			expectedErr: errorx.BadRequestError("Bad Request"),
		},
		{
			name:        "Unauthorized",
			fiberError:  fiber.NewError(fiber.StatusUnauthorized, "Unauthorized"),
			expectedErr: errorx.UnauthorizedError("Unauthorized"),
		},
		{
			name:        "Forbidden",
			fiberError:  fiber.NewError(fiber.StatusForbidden, "Forbidden"),
			expectedErr: errorx.ForbiddenError("Forbidden"),
		},
		{
			name:        "NotFound",
			fiberError:  fiber.NewError(fiber.StatusNotFound, "Not Found"),
			expectedErr: errorx.NotFoundError("Not Found"),
		},
		{
			name:        "MethodNotAllowed",
			fiberError:  fiber.NewError(fiber.StatusMethodNotAllowed, "Method Not Allowed"),
			expectedErr: errorx.MethodNotAllowedError("Method Not Allowed"),
		},
		{
			name:        "RequestTimeout",
			fiberError:  fiber.NewError(fiber.StatusRequestTimeout, "Request Timeout"),
			expectedErr: errorx.TimeoutError("Request Timeout"),
		},
		{
			name:        "Conflict",
			fiberError:  fiber.NewError(fiber.StatusConflict, "Conflict"),
			expectedErr: errorx.ConflictError("Conflict"),
		},
		{
			name:        "TooManyRequests",
			fiberError:  fiber.NewError(fiber.StatusTooManyRequests, "Too Many Requests"),
			expectedErr: errorx.TooManyRequestError("Too Many Requests"),
		},
		{
			name:        "InternalServerError",
			fiberError:  fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error"),
			expectedErr: errorx.InternalServerError("Internal Server Error"),
		},
		{
			name:        "UnknownErrorCode",
			fiberError:  fiber.NewError(999, "Custom Error"),
			expectedErr: errorx.Failed("Custom Error"),
		},
		{
			name:        "NilFiberError",
			fiberError:  nil,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapFiberError(tt.fiberError)

			if tt.expectedErr == nil {
				assert.Nil(t, result)
			} else {
				e, ok := result.(*errorx.Error)
				expectedErr := tt.expectedErr.(*errorx.Error)

				assert.True(t, ok)
				assert.Equal(t, expectedErr.Code, e.Code)
				assert.Equal(t, expectedErr.Message, e.Message)
				assert.Equal(t, expectedErr.Status, e.Status)
			}
		})
	}
}

type Struct struct {
	Data struct {
		Name string
	} `json:"data"`
}

func TestJson(t *testing.T) {
	jsonData := `{"data": ["hello"]}`
	var result Struct
	err := json.Unmarshal([]byte(jsonData), &result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	t.Log(result)

}
