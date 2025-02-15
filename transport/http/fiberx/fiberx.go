package fiberx

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/rewrite"
	"github.com/gofiber/swagger"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/metadata"
	"github.com/kingstonduy/go-core/metrics"
	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/transport"
	"github.com/kingstonduy/go-core/validation"
)

type FiberAppOptions struct {
	ServiceName string

	ServiceVersion string

	Validator validation.Validator

	Tracer trace.Tracer

	Logger logger.Logger

	Metrics *metrics.Metrics

	HealthCheckEnabled bool

	SwaggerEnabled bool

	LoggingRequestEnabled bool

	RateLimiterEnabled bool

	RequestTracingEnabled bool

	MetricEndpointEnabled bool

	BaseRequestBodyValidationEnabled bool

	APIRegistrationLoggingEnabled bool

	SwaggerPath string

	MetricsPath string

	LivenessPath string

	LivenessProbe func(c *fiber.Ctx) bool

	ReadinessPath string

	ReadinessProbe func(c *fiber.Ctx) bool

	BasePath string

	RateLimiterConfig RateLimiterConfig

	RecoverHandlerEnabled bool

	// fiber configurations
	FiberConfig *fiber.Config
}

type FiberAppOption func(*FiberAppOptions)

func NewFiberOptions(opts ...FiberAppOption) FiberAppOptions {
	// default options
	options := FiberAppOptions{
		ServiceName:                      metadata.DefaultServiceName,
		ServiceVersion:                   metadata.DefaultServiceName,
		HealthCheckEnabled:               true,
		SwaggerEnabled:                   true,
		RateLimiterEnabled:               true,
		LoggingRequestEnabled:            true,
		RequestTracingEnabled:            true,
		MetricEndpointEnabled:            true,
		BaseRequestBodyValidationEnabled: true,
		APIRegistrationLoggingEnabled:    true,
		SwaggerPath:                      FiberAppSwaggerEndpoint,
		MetricsPath:                      FiberAppMetricsEndpoint,
		RateLimiterConfig:                DefaultRateLimiterConfig,
		RecoverHandlerEnabled:            true,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

// Service name
// Default: metadata.DefaultServiceName
func WithServiceName(serviceName string) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.ServiceName = serviceName
	}
}

// Service version
// Default: metadata.DefaultServiceVersion
func WithServiceVersion(serviceVersion string) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.ServiceVersion = serviceVersion
	}
}

// Validator used in fiber app for base request validation
// Default: validation.DefaultValidator
func WithValidator(validator validation.Validator) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.Validator = validator
	}
}

// Tracer used in fiber app for request monitoring
// Default: trace.DefaultTracer
func WithTracer(tracer trace.Tracer) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.Tracer = tracer
	}
}

// Logger used in fiber app for request monitoring, loggings
// Default: logger.DefaultLogger
func WithLogger(logger logger.Logger) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.Logger = logger
	}
}

// Metrics used in fiber app for request metrics if needed
// Default: metrics.Default()
func WithMetrics(metrics *metrics.Metrics) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.Metrics = metrics
	}
}

// Enable health check endpoints (liveness, readiness)
// Default: true
func WithHealthCheckEnabled(enabled bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.HealthCheckEnabled = enabled
	}
}

// Enable swagger
// Default: true
func WithSwaggerEnabled(enabled bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.SwaggerEnabled = enabled
	}
}

// Enable logging request status code, route, method
// Default: true
func WithLoggingRequestEnabled(enabled bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.LoggingRequestEnabled = enabled
	}
}

// Enable request rate limit
// Default: true
func WithRateLimiterEnabled(enabled bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.RateLimiterEnabled = enabled
	}
}

// Opentelemetry tracing, to inject traceID from request headers or create new traceID for each request
// Trace request duration and information in the trace segment on each request
// Default: true
func WithRequestTracingEnabled(enabled bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.RequestTracingEnabled = enabled
	}
}

// Metrics endpoint enabled for HTTP
// Default: true
func WithMetricEndpointEnabled(enabled bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.MetricEndpointEnabled = enabled
	}
}

// Validate base request body enabled
// Default: true
func WithBaseRequestBodyValidationEnabled(enabled bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.BaseRequestBodyValidationEnabled = enabled
	}
}

// Print register API
// Default: true
func WithAPIRegistrationLoggingEnabled(enabled bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.APIRegistrationLoggingEnabled = enabled
	}
}

// Swagger endpoint
// Default: /swagger/*
func WithSwaggerPath(path string) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.SwaggerPath = path
	}
}

// Metrics endpoint
// Default: /metrics
func WithMetricsPath(path string) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.MetricsPath = path
	}
}

// Liveness endpoint
// Default: /livez
func WithLivenessPath(path string) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.LivenessPath = path
	}
}

// Checks if the server is up and running.
// Default: By default returns true immediately when the server is operational.
func WithLivenessProbe(probe func(c *fiber.Ctx) bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.LivenessProbe = probe
	}
}

// Readiness endpoint
// Default: /readyz
func WithReadinessPath(path string) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.ReadinessPath = path
	}
}

// Assesses if the application is ready to handle requests.
// Default: By default returns true immediately when the server is operational.
func WithReadinessProbe(probe func(c *fiber.Ctx) bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.ReadinessProbe = probe
	}
}

// Base Path
// Default: ""
func WithBasePath(basePath string) FiberAppOption {
	return func(options *FiberAppOptions) {
		basePath := strings.TrimSuffix(basePath, "/")
		options.BasePath = basePath
	}
}

func WithRateLimiterConfig(config RateLimiterConfig) FiberAppOption {
	return func(options *FiberAppOptions) {
		// set default value
		if config.Max <= 0 {
			config.Max = DefaultRateLimiterConfig.Max
		}

		if int(config.Duration.Seconds()) <= 0 {
			config.Duration = DefaultRateLimiterConfig.Duration
		}

		if config.LimitReached == nil {
			config.LimitReached = DefaultRateLimiterConfig.LimitReached
		}

		options.RateLimiterConfig = config
	}
}

func WithRecoverHandlerEnabled(enabled bool) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.RecoverHandlerEnabled = enabled
	}
}

func WithFiberConfig(config fiber.Config) FiberAppOption {
	return func(options *FiberAppOptions) {
		options.FiberConfig = &config
	}
}

const (
	FiberAppSwaggerEndpoint = "/swagger/*"
	FiberAppMetricsEndpoint = "/metrics"
	StepMonitoringRequest   = "request-monitoring"
	FiberAppMonitorEndpoint = "/monitor"
)

type FiberApp struct {
	*fiber.App
	Options FiberAppOptions
}

func NewFiberApp(opts ...FiberAppOption) *FiberApp {
	options := NewFiberOptions(opts...)

	fiberConfig := DefaultFiberConfig

	if optFiberConfig := options.FiberConfig; optFiberConfig != nil {
		fiberConfig = *optFiberConfig
	}

	fiberApp := fiber.New(fiberConfig)

	app := FiberApp{
		App:     fiberApp,
		Options: options,
	}

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,                                                  // Allow cookies and other credentials
		AllowOrigins:     "http://localhost:5173,https://kingstonduy.github.io", // Specify the allowed origin
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",              // Specify allowed methods
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	// Handle global panic
	if app.Options.RecoverHandlerEnabled {
		app.Use(NewRecoverHandler())
	}

	// Base Path
	if basePath := app.Options.BasePath; basePath != "" {
		app.Use(rewrite.New(rewrite.Config{
			Rules: map[string]string{
				fmt.Sprintf("%s/*", basePath): "/$1",
			},
		}))
	}

	// HeathCheck
	if app.Options.HealthCheckEnabled {
		healthcheckConfig := healthcheck.Config{
			LivenessEndpoint:  app.Options.LivenessPath,
			LivenessProbe:     app.Options.LivenessProbe,
			ReadinessEndpoint: app.Options.ReadinessPath,
			ReadinessProbe:    app.Options.ReadinessProbe,
		}
		app.Use(healthcheck.New(healthcheckConfig))
	}

	// Logging HTTP
	if app.Options.LoggingRequestEnabled {
		fiberApp.Use(fiberLog.New(fiberLog.Config{
			Next:         nil,
			Done:         nil,
			Format:       "[${time}] ${status} ${latency} ${method} ${path}\n",
			TimeFormat:   "2006-01-02 15:04:05",
			TimeZone:     "Local",
			TimeInterval: 500 * time.Millisecond,
			Output:       os.Stdout,
		}))
	}

	// Swagger
	if options.SwaggerEnabled {
		swaggerConfig := swagger.Config{
			URL: "doc.json",
		}

		if srvName := app.Options.ServiceName; len(srvName) > 0 {
			swaggerConfig.Title = srvName
		}

		fiberApp.Get(options.SwaggerPath, swagger.New(swaggerConfig))
	}

	// Prometheus metrics
	if app.Options.MetricEndpointEnabled {
		prometheus := fiberprometheus.New(app.Options.ServiceName)
		prometheus.RegisterAt(fiberApp, options.MetricsPath)
		fiberApp.Use(prometheus.Middleware)
	}

	// request tracing and monitoring
	if app.Options.RequestTracingEnabled {
		// Inject or create TraceID for each request
		fiberApp.Use(otelfiber.Middleware())

		// Request monitoring: duration, information in the request body
		fiberApp.Use(func(ctx *fiber.Ctx) error {
			var req transport.Request[any]
			ctx.BodyParser(&req) // nolint
			trace := req.Trace

			reqHeaders := make(map[string]string)
			// extract request headers
			for k, v := range ctx.GetReqHeaders() {
				reqHeaders[k] = strings.Join(v, ", ")
			}

			monitoredCtx := transport.MonitorRequest(
				ctx.UserContext(),
				transport.MonitorRequestData{
					ClientIP:           ctx.IP(),
					Protocol:           metadata.ProtocolHTTP,
					Method:             ctx.Method(),
					RequestPath:        app.getPath(ctx.Path(), true),
					Hostname:           ctx.Hostname(),
					UserAgent:          string(ctx.Context().UserAgent()),
					ClientTime:         trace.Cts,
					ContentLength:      len(ctx.Request().Body()),
					Request:            req,
					ClientID:           trace.Cid,
					From:               trace.From,
					To:                 trace.To,
					Username:           trace.Username,
					RequestHeaders:     reqHeaders,
					RemoteHost:         ctx.Context().RemoteAddr().String(),
					XForwardedFor:      strings.Join(ctx.IPs(), ","),
					TransactionTimeout: trace.TransactionTimeout,
					SystemID:           trace.Sid,
				},
				transport.WithLogger(app.getLogger()),
				transport.WithTracer(app.getTracer()),
			)

			// Inject the context for later
			ctx.SetUserContext(monitoredCtx)

			return ctx.Next()
		})
	}

	// Validate base request format
	// if app.Options.BaseRequestBodyValidationEnabled {
	// 	fiberApp.Use(func(ctx *fiber.Ctx) error {
	// 		var req transport.Request[any]
	// 		err := ctx.BodyParser(&req)
	// 		if err != nil {
	// 			return errorx.BadRequestError("Failed to parse request: Invalid base request format. %v", err)
	// 		}

	// 		err = app.getValidator().Validate(req)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		return ctx.Next()
	// 	})
	// }

	// Hooks
	// Register route hook to print routes registered with fiber app
	if app.Options.APIRegistrationLoggingEnabled {
		fiberApp.Hooks().OnRoute(func(route fiber.Route) error {
			app.getLogger().Infof(context.Background(), "[HTTP] Registered route: %s %s", route.Method, fmt.Sprintf("%s%s", app.Options.BasePath, route.Path))
			return nil
		})
	}

	// Rate Limiter
	if options.RateLimiterEnabled {
		config := options.RateLimiterConfig
		limiterConfig := limiter.Config{
			Max:                    config.Max,
			Expiration:             config.Duration,
			LimitReached:           config.LimitReached,
			SkipFailedRequests:     config.SkipFailedRequests,
			SkipSuccessfulRequests: config.SkipSuccessfulRequests,
		}

		app.Use(limiter.New(limiterConfig))
	}

	return &app
}

func (app *FiberApp) getLogger() logger.Logger {
	if app.Options.Logger != nil {
		return app.Options.Logger
	}
	return logger.DefaultLogger
}

func (app *FiberApp) getTracer() trace.Tracer {
	if app.Options.Tracer != nil {
		return app.Options.Tracer
	}
	return trace.DefaultTracer
}

func (app *FiberApp) getValidator() validation.Validator {
	if app.Options.Validator != nil {
		return app.Options.Validator
	}
	return validation.DefaultValidator
}

// Get the path with the given
func (app *FiberApp) getPath(path string, includeBasePath bool) string {
	if !includeBasePath {
		return path
	}

	basePath := app.Options.BasePath
	return fmt.Sprintf("%s%s", basePath, path)
}
