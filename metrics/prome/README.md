#### Usage

```go
func GetPrometheusMetricer(logger logger.Logger) *prome.PrometheusMetricer {
	metricer, err := prome.NewPrometheusMetricer()

	if err != nil {
		logger.Error(context.TODO(), "Failed to create prometheous metricer")
		panic(err)
	}

	return metricer
}

```

```go

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
	reqType := util.GetType(request)

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		b.metricer.RequestSummary.WithLabelValues(b.cfg.Name, reqType).Observe(us)
		b.metricer.RequestHistogram.WithLabelValues(b.cfg.Name, reqType).Observe(v)
	}))

	defer func() {
		var businessErr *domainErrors.DomainError
		// mark a request success as if there is no error happened or the error is business error
		if err == nil || errors.As(err, &businessErr) {
			b.metricer.RequestTotalCounter.WithLabelValues(b.cfg.Name, reqType, "success").Inc()
		} else {
			b.metricer.RequestTotalCounter.WithLabelValues(b.cfg.Name, reqType, "failure").Inc()
		}

		// stop timer
		timer.ObserveDuration()
	}()

	response, err = next(ctx)
	return response, err
}
```