package fiberx

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RateLimiterConfig defines the config for middleware.
type RateLimiterConfig struct {
	// Max number of recent connections during `Duration` seconds before sending a 429 response
	//
	// Default: 500
	Max int

	// Duration is the time on how long to keep records of requests in memory
	//
	// Default: 1 * time.Second
	Duration time.Duration

	// LimitReached is called when a request hits the limit
	//
	// Default: func(c *fiber.Ctx) error {
	//   return fiber.ErrTooManyRequests
	// }
	LimitReached fiber.Handler

	// When set to true, requests with StatusCode >= 400 won't be counted.
	//
	// Default: false
	SkipFailedRequests bool

	// When set to true, requests with StatusCode < 400 won't be counted.
	//
	// Default: false
	SkipSuccessfulRequests bool
}

var DefaultRateLimiterConfig = RateLimiterConfig{
	Max:      500,
	Duration: 1 * time.Second,
	LimitReached: func(c *fiber.Ctx) error {
		return fiber.ErrTooManyRequests
	},
	SkipFailedRequests:     false,
	SkipSuccessfulRequests: false,
}

var DefaultFiberConfig = fiber.Config{
	ErrorHandler:          CustomErrorHandler,
	JSONDecoder:           json.Unmarshal,
	JSONEncoder:           json.Marshal,
	DisableStartupMessage: true,
}
