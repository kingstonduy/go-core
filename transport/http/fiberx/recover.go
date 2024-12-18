package fiberx

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/logger"
)

func DefaultStackTraceHandler(c *fiber.Ctx, e interface{}) {
	logger.Errorf(c.UserContext(), "panic: %v\n%s\n", e, debug.Stack()) //nolint:errcheck // This will never fail
}

// RecoverHandler: handle global panic to prevent crash app when panic
func NewRecoverHandler() fiber.Handler {
	// Return new handler
	return func(c *fiber.Ctx) (err error) { //nolint:nonamedreturns // Uses recover() to overwrite the error
		// Catch panics
		defer func(ctx context.Context) {
			if r := recover(); r != nil {
				// print stack trace
				DefaultStackTraceHandler(c, r)

				// return internal server error
				err = errorx.InternalServerError(fmt.Sprintf("%v", r))
			}
		}(c.UserContext())

		// Return err if exist, else move to next handler
		return c.Next()
	}
}
