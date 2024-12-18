#### Usage

```go

// USECASE
type CheckBalanceRequest struct {
	Account string `json:"account"`
}

type CheckBalanceResponse struct {
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type CheckBalanceHandler interface {
	Handle(ctx context.Context, request *CheckBalanceRequest) (*CheckBalanceResponse, error)
}

// Register handlers
func RegisterPipelineHandlers(
	checkBalanceHandler account.CheckBalanceHandler,
	// other handers
) {
	pipeline.RegisterRequestHandler(checkBalanceHandler)
	// register other handlers
}

// Register behaviours
func RegisterPipelineBehaviors(
	requestLoggingBehavior *RequestLoggingBehavior,
	requestTracingBehavior *RequestTracingBehavior,
	requestMetricBehavior *RequestMetricBehavior,
	errorHandlingBehavior *ErrorHandlingBehavior,
) {
	pipeline.RegisterRequestPipelineBehaviors(
		requestTracingBehavior,
		requestLoggingBehavior,
		requestMetricBehavior,
		errorHandlingBehavior,
	)
}

// ERROR HANDLING FOR RECOVERING FROM PANIC
type ErrorHandlingBehavior struct {
	logger logger.Logger
}

func NewErrorHandlingBehavior(logger logger.Logger) *ErrorHandlingBehavior {
	return &ErrorHandlingBehavior{
		logger: logger,
	}
}

func (b *ErrorHandlingBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (res interface{}, err error) {
	response, err := next(ctx)
	return response, err
}

// TRACING
type RequestTracingBehavior struct {
	logger logger.Logger
	tracer trace.Tracer
}

func NewTracingBehavior(logger logger.Logger, tracer trace.Tracer) *RequestTracingBehavior {
	return &RequestTracingBehavior{
		logger: logger,
		tracer: tracer,
	}
}

func (b *RequestTracingBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (interface{}, error) {
	res, err := next(ctx)
	return res, err
}

// METRICS
type RequestMetricBehavior struct {
	logger   logger.Logger
	metricer *prome.PrometheusMetricer
	cfg      *ServerConfig
}

func NewMetricBehavior(logger logger.Logger, metricer *prome.PrometheusMetricer, cfg *ServerConfig) *RequestMetricBehavior {
	return &RequestMetricBehavior{
		logger:   logger,
		metricer: metricer,
		cfg:      cfg,
	}
}

func (b *RequestMetricBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (response interface{}, err error) {
	response, err = next(ctx)
	return response, err
}

// LOGGING
type RequestLoggingBehavior struct {
	logger logger.Logger
}

func NewRequestLoggingBehavior(logger logger.Logger) *RequestLoggingBehavior {
	return &RequestLoggingBehavior{
		logger: logger,
	}
}

func (b *RequestLoggingBehavior) Handle(ctx context.Context, request interface{}, next pipeline.RequestHandlerFunc) (response interface{}, err error) {
	response, err = next(ctx)
	return response, err
}

```