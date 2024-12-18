package hresty

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/util"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type RestyOptions struct {
	LoggingRequestEnable  bool
	LoggingResponseEnable bool
	LoggingErrorEnable    bool
	Timeout               *time.Duration
	TlsConfig             *tls.Config
}

type RestyOption func(*RestyOptions)

func NewRestyClient(opts ...RestyOption) *resty.Client {
	client := resty.New()
	options := NewRestyOptions(opts...)

	// Set timeout
	if options.Timeout != nil {
		client.SetTimeout(*options.Timeout)
	}
	// tracing instrumentation
	transport := http.DefaultTransport.(*http.Transport)
	// Set TSL configuration
	if options.TlsConfig != nil {
		transport.TLSClientConfig = options.TlsConfig
	}
	tracedTransport := otelhttp.NewTransport(transport)
	client.SetTransport(tracedTransport)

	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		ctx := r.Context()
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

		if options.LoggingRequestEnable {
			method := strings.ToUpper(r.Method)
			stepName := fmt.Sprintf("Request Method %s - RESTY", method)
			// Log the request
			logger.Fields(map[string]interface{}{
				logger.FIELD_OPERATOR_NAME: r.URL,
				logger.FIELD_STEP_NAME:     stepName,
			}).Infof(ctx, util.MakeStringLogs(ctx, r.Header, r.Body, ""))

		}

		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		ctx := r.Request.Context()
		if options.LoggingResponseEnable {
			method := strings.ToUpper(r.Request.Method)
			stepName := fmt.Sprintf("Response Method %s - RESTY", method)
			// Log the response
			logger.Fields(map[string]interface{}{
				logger.FIELD_OPERATOR_NAME: r.Request.URL,
				logger.FIELD_STEP_NAME:     stepName,
			}).Infof(ctx, util.MakeStringLogs(ctx, r.Header(), r.Body(), r.Status()))
		}
		return nil
	})

	client.OnError(func(r *resty.Request, err error) {
		ctx := r.Context()
		if options.LoggingErrorEnable {
			method := strings.ToUpper(r.Method)
			stepName := fmt.Sprintf("Error Method %s - RESTY", method)
			// Log the Error
			logger.Fields(map[string]interface{}{
				logger.FIELD_OPERATOR_NAME: r.URL,
				logger.FIELD_STEP_NAME:     stepName,
			}).Errorf(ctx, err.Error())
		}
	})

	return client
}

func NewRestyOptions(opts ...RestyOption) RestyOptions {
	// default options
	options := RestyOptions{
		LoggingRequestEnable:  true,
		LoggingResponseEnable: true,
		LoggingErrorEnable:    true,
	}

	for _, opt := range opts {
		opt(&options)
	}
	return options
}

func WithLoggingRequestEnable(enable bool) RestyOption {
	return func(options *RestyOptions) {
		options.LoggingRequestEnable = enable
	}
}

func WithLoggingResponseEnable(enable bool) RestyOption {
	return func(options *RestyOptions) {
		options.LoggingResponseEnable = enable
	}
}

func WithLoggingErrorEnable(enable bool) RestyOption {
	return func(options *RestyOptions) {
		options.LoggingErrorEnable = enable
	}
}

func WithTimeOut(timeout *time.Duration) RestyOption {
	return func(options *RestyOptions) {
		options.Timeout = timeout
	}
}

func WithTlsConfig(tlsConfig *tls.Config) RestyOption {
	return func(options *RestyOptions) {
		options.TlsConfig = tlsConfig
	}
}
