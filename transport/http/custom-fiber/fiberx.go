package Fiberx

import (
	"encoding/json"

	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
	errorx "github.com/kingstonduy/go-core/error"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/register"
	transport "github.com/kingstonduy/go-core/transport/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Fiberx struct {
	App                      *fiber.App
	isSwaggerEnabled         bool
	isMonitorEnabled         bool
	isValidateRequestEnabled bool
	isTraceEnabled           bool
}

type option func(f *Fiberx)

func NewFiberx(opts ...option) Fiberx {
	// default options
	fiberx := Fiberx{
		App:                      fiber.New(),
		isSwaggerEnabled:         true,
		isMonitorEnabled:         true,
		isValidateRequestEnabled: true,
		isTraceEnabled:           true,
	}

	for _, opt := range opts {
		opt(&fiberx)
	}

	if fiberx.isMonitorEnabled {
		fiberx.App.Use("/metrics", monitor.New())
	}

	if fiberx.isSwaggerEnabled {
		swaggerConfig := swagger.Config{URL: "doc.json"}
		fiberx.App.Use("/swagger/*", swagger.New(swaggerConfig))
	}

	if fiberx.isTraceEnabled {
		fiberx.App.Use(otelfiber.Middleware(otelfiber.WithPropagators(otel.GetTextMapPropagator())))
	}

	return fiberx
}

// T and K should be transport.RequestType and transport.ResponseType
func FiberPipeline[T any, K any](Fiberx Fiberx, c *fiber.Ctx, condtion string) {
	// tracing
	Fiberx.App.Use(otelfiber.Middleware())

	ctx := c.UserContext()
	var t transport.RequestType[T]
	var res transport.ResponseType[K]

	// 6. response result to client
	defer func() {
		// set response headers
		headers := c.GetRespHeaders()

		// Inject OpenTelemetry trace information into the headers
		otel.GetTextMapPropagator().Inject(c.UserContext(), propagation.HeaderCarrier(headers))

		reqHeaders := make(map[string]string)
		// extract request headers
		for k, v := range headers {
			reqHeaders[k] = strings.Join(v, ", ")
		}
		for k, v := range reqHeaders {
			c.Set(k, v)
		}

		// set response body
		c.Status(res.Result.Status).JSON(res)
	}()

	// 1. parse json
	if err := json.Unmarshal(c.Body(), &t); err != nil {
		logger.Error(ctx, err)
		errx := errorx.InvalidData(errorx.WithDetail(err.Error()))
		res.Result = transport.GetResultFromErrorx(*errx)
		return
	}

	// 2. logging request
	reqBytes, _ := json.Marshal(t)
	logger.Infof(ctx, "request: %s", string(reqBytes))

	// 5. logging response
	defer func() {
		resBytes, _ := json.Marshal(res)
		logger.Infof(ctx, "response: %s", string(resBytes))
	}()

	// 3. validate request
	validate := validator.New()
	err := validate.Struct(t)
	if err != nil {
		logger.Error(ctx, err)
		errx := errorx.InvalidData(errorx.WithDetail(err.Error()))
		res.Result = transport.GetResultFromErrorx(*errx)
		return
	}

	// 4. process request
	res, err = register.Process[transport.RequestType[T], transport.ResponseType[K]](ctx, condtion, t)
	res.Trace = transport.GetResponseTrace(t.Trace)
	if err != nil { // wrapper
		if errx, ok := err.(*errorx.Errorx); ok {
			res.Result = transport.GetResultFromErrorx(*errx)
			return
		} else { //return internal server error
			errx := errorx.InternalServerErrorx(errorx.WithDetail(err.Error()))
			res.Result = transport.GetResultFromErrorx(*errx)
		}
	} else {
		// no error at all'
		res.Result = transport.GetSuccessResult()
		return
	}
}
