package fiberx

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/pipeline"
	"github.com/kingstonduy/go-core/transport"
)

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		err = wrapFiberError(fiberError)
	}

	fRes := transport.GetResponse[interface{}](
		ctx.UserContext(),
		transport.WithError(err),
	)

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return ctx.Status(fRes.Result.StatusCode).JSON(fRes)
}

// Handle request for fiber
// error: system error, not API error
func RequestHandler[TReq any, TResp any](ctx *fiber.Ctx, opts ...RequestHandlerOption) error {
	options := NewRequestHandlerOptions(opts...)

	// Step 1: Parse the request
	var req transport.Request[TReq]
	err := ctx.BodyParser(&req)
	if err != nil {
		return errorx.BadRequestError("Failed to parse request: Invalid base request format. %v", err)
	}

	// Create empty object if the req.Data is nil
	if reflect.ValueOf(req.Data).Kind() == reflect.Ptr && reflect.ValueOf(req.Data).IsNil() {
		t := reflect.TypeOf(req.Data)
		// Create a new instance of the underlying type, as a pointer
		newInstance := reflect.New(t.Elem()).Interface()
		// Assign new pointer instance to req.Data
		req.Data = newInstance.(TReq)
	}

	// Step 2: Send request to the request pipeline
	resp, err := pipeline.Send[TReq, TResp](ctx.UserContext(), req.Data)

	// Step 3: Parse response
	httpResp := transport.GetResponse[TResp](
		ctx.UserContext(),
		transport.WithData(resp),
		transport.WithError(err),
	)

	// trigger hook when request handled
	if options.OnRequestHandledFunc != nil {
		go func() {
			options.OnRequestHandledFunc(ctx.UserContext(), httpResp)
		}()
	}

	// Step 4: Send HTTP Response
	return ctx.Status(httpResp.Result.StatusCode).JSON(httpResp)
}

type RequestHandlerOptions struct {
	OnRequestHandledFunc func(ctx context.Context, res interface{})
	PreHandlerEnable     bool
	PreHandlerFunc       func(ctx *fiber.Ctx)
	PostHandlerFunc      func(ctx *fiber.Ctx, b []byte)
	PostHandlerEnable    bool
	VerifyTokenEnable    bool
}

type RequestHandlerOption func(*RequestHandlerOptions)

func NewRequestHandlerOptions(opts ...RequestHandlerOption) *RequestHandlerOptions {
	options := RequestHandlerOptions{}

	for _, opt := range opts {
		opt(&options)
	}

	return &options
}

func WithPreHandlerFunc(f func(ctx *fiber.Ctx)) RequestHandlerOption {
	return func(opts *RequestHandlerOptions) {
		opts.PreHandlerEnable = true
		opts.PreHandlerFunc = f
	}
}

func WithPostHandlerFunc(f func(ctx *fiber.Ctx, b []byte)) RequestHandlerOption {
	return func(opts *RequestHandlerOptions) {
		opts.PostHandlerEnable = true
		opts.PostHandlerFunc = f
	}
}

func WithAuthentication() RequestHandlerOption {
	return func(opts *RequestHandlerOptions) {
		// TODO enable it = true
		opts.VerifyTokenEnable = false
	}
}

// Handle after the request handled. res type: transport.Response[T]
func WithOnRequestHandledFunc(f func(ctx context.Context, res interface{})) RequestHandlerOption {
	return func(opts *RequestHandlerOptions) {
		opts.OnRequestHandledFunc = f
	}
}

func RequestHandlerWithDynamicTimeout[TReq any, TResp any](ctx *fiber.Ctx, opts ...RequestHandlerOption) error {
	options := NewRequestHandlerOptions(opts...)

	if options.PreHandlerEnable {
		options.PreHandlerFunc(ctx)
	}

	if options.VerifyTokenEnable {
		token := ctx.Cookies("jwt", "")
		if token == "" {
			httpResp := transport.GetResponse[TResp](
				ctx.UserContext(),
				transport.WithError(errorx.AuthenticationErrorWithDetails("Missing token", "")),
			)
			return ctx.Status(httpResp.Result.StatusCode).JSON(httpResp)
		}

		// Verify token here
		if _, err := VerifyToken(token); err != nil {
			httpResp := transport.GetResponse[TResp](
				ctx.UserContext(),
				transport.WithError(errorx.AuthenticationErrorWithDetails(err.Error(), "")),
			)
			return ctx.Status(httpResp.Result.StatusCode).JSON(httpResp)
		}
	}
	var req transport.Request[TReq]
	err := ctx.BodyParser(&req)
	if err != nil {
		return errorx.BadRequestError("Failed to parse request: Invalid base request format. %v", err)
	}
	// Define timeout duration default
	timeoutDuration := time.Duration(60000) * time.Millisecond
	if req.Trace.TransactionTimeout > 0 {
		timeoutDuration = time.Duration(req.Trace.TransactionTimeout) * time.Millisecond
	}
	logger.Info(ctx.UserContext(), "Time duration for request timeout: ", timeoutDuration)
	// Create a context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx.UserContext(), timeoutDuration)
	defer cancel()

	// Step 1: Parse the request
	// Step 2: Send request to the request pipeline with the contextWithTimeout
	respChan := make(chan *transport.Response[TResp], 1)
	go func() {
		resp, err := pipeline.Send[TReq, TResp](ctxWithTimeout, req.Data)
		httpResp := transport.GetResponse[TResp](
			ctxWithTimeout,
			transport.WithData(resp),
			transport.WithError(err),
		)
		respChan <- &httpResp
	}()

	select {
	case <-ctxWithTimeout.Done():
		err := errorx.TimeoutErrorWithDetails(errorx.ErrorMessageTimeout, "")
		httpResp := transport.GetResponse[TResp](
			ctxWithTimeout,
			transport.WithError(err),
		)
		logger.Infof(ctx.Context(), "%s", err.Error())
		return ctx.Status(httpResp.Result.StatusCode).JSON(httpResp)

	case resp := <-respChan:

		// Trigger hook when request handled
		if options.OnRequestHandledFunc != nil {
			options.OnRequestHandledFunc(ctx.UserContext(), resp)
		}

		if options.PostHandlerEnable {
			b, _ := json.Marshal(resp)
			options.PostHandlerFunc(ctx, b)
		}
		// Step 4: Send HTTP Response
		return ctx.Status(resp.Result.StatusCode).JSON(resp)
	}
}

func wrapFiberError(fError *fiber.Error) error {
	if fError == nil {
		return nil
	}

	switch fError.Code {
	case fiber.StatusBadRequest: // 400
		return errorx.BadRequestError(fError.Message)
	case fiber.StatusUnauthorized: // 401
		return errorx.UnauthorizedError(fError.Message)
	case fiber.StatusForbidden: // 403
		return errorx.ForbiddenError(fError.Message)
	case fiber.StatusNotFound: // 404
		return errorx.NotFoundError(fError.Message)
	case fiber.StatusMethodNotAllowed: // 405
		return errorx.MethodNotAllowedError(fError.Message)
	case fiber.StatusRequestTimeout: // 408
		return errorx.TimeoutError(fError.Message)
	case fiber.StatusConflict: // 409
		return errorx.ConflictError(fError.Message)
	case fiber.StatusTooManyRequests: // 429
		return errorx.TooManyRequestError(fError.Message)
	case fiber.StatusInternalServerError: // 500
		return errorx.InternalServerError(fError.Message)
	default:
		return errorx.Failed(fError.Message)
	}
}
