package pipeline

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/metadata"
	"github.com/kingstonduy/go-core/metrics"
	"github.com/kingstonduy/go-core/metrics/prometheus"
	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/trace/otel"
	"github.com/kingstonduy/go-core/validation"
	"github.com/kingstonduy/go-core/validation/goplayaround"
)

func GetTracer(t *testing.T) trace.Tracer {
	tracer, err := otel.NewOpenTelemetryTracer(
		context.Background(),
		trace.WithTraceServiceName("test_behavior_service"),
		trace.WithServiceVersion("1.0.0"),
		trace.WithTraceExporterEndpoint("localhost:4318"),
	)

	if err != nil {
		t.Fatalf("ErrorL %v", err)
	}
	return tracer
}

func GetMetrics(t *testing.T) *metrics.Metrics {
	promeConfigs := prometheus.PrometheusOpts{
		Expiration: 0,
		Name:       "prometheus_metrics",
	}

	sink, err := prometheus.NewPrometheusSinkFrom(promeConfigs)
	if err != nil {
		t.Fatalf("ErrorL %v", err)
	}

	m, err := metrics.New(metrics.DefaultConfig(metadata.DefaultServiceName), sink)
	if err != nil {
		t.Fatalf("ErrorL %v", err)
	}

	return m

}

func GetValidator(t *testing.T) validation.Validator {
	return goplayaround.NewGpValidator()
}

func GetLogger(t *testing.T) logger.Logger {
	return logger.DefaultLogger
}

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

type Request struct {
	Number float64 `validate:"min=100"`
}

type Response struct {
	Result float64
}

func (h *handler) Handle(ctx context.Context, req *Request) (*Response, error) {
	time.Sleep(100 * time.Millisecond)
	if req.Number == 0 {
		panic("System failed")
	}

	if req.Number == 1 {
		return nil, fmt.Errorf("Error: number is one")
	}

	return &Response{
		Result: math.Pow(req.Number, 2),
	}, nil
}

func TestPipelineBehaviorsWithDefaultValues(t *testing.T) {
	handler := NewHandler()
	if err := RegisterRequestHandler(handler); err != nil {
		t.Error(err)
	}

	_, err := Send[*Request, *Response](context.Background(), &Request{
		Number: 999,
	})

	if err != nil {
		t.Logf("%v", err.Error())
	} else {
		t.Log("OK")
	}

	time.Sleep(5 * time.Second)

}

func TestPipelineBehaviors(t *testing.T) {
	handler := NewHandler()
	if err := RegisterRequestHandler(handler); err != nil {
		t.Error(err)
	}

	_, err := Send[*Request, *Response](context.Background(), &Request{
		Number: 999,
	})

	if err != nil {
		t.Logf("%v", err.Error())
	} else {
		t.Log("OK")
	}

	time.Sleep(5 * time.Second)

}

func TestPipelineBehaviorsWithDefaultValuesWasSetByDefault(t *testing.T) {
	handler := NewHandler()
	if err := RegisterRequestHandler(handler); err != nil {
		t.Error(err)
	}

	trace.SetDefaultTracer(GetTracer(t))
	logger.SetDefaultLogger(GetLogger(t))
	validation.SetDefaultValidator(GetValidator(t))
	if _, err := metrics.NewGlobal(metrics.DefaultConfig(metadata.DefaultServiceName), &prometheus.PrometheusSink{}); err != nil {
		t.Error(err)
	}

	_, err := Send[*Request, *Response](context.Background(), &Request{
		Number: 999,
	})

	if err != nil {
		t.Logf("%v", err.Error())
	} else {
		t.Log("OK")
	}

	time.Sleep(5 * time.Second)

}
