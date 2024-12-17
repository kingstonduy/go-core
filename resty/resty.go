package resty

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/kingstonduy/go-core/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type RestyOptions struct {
	LoggingRequestEnable  bool
	LoggingResponseEnable bool
}

func MakeStringLogs(ctx context.Context, headers map[string][]string, body interface{}, httpStatus string) string {
	object := make(map[string]interface{})

	// Add headers
	if headers != nil {
		object["headers"] = headers
	} else {
		object["headers"] = "empty"
	}

	// Add HTTP status if not empty
	if httpStatus != "" {
		object["http_status"] = httpStatus
	}

	// Add body
	if body != nil {
		if _, ok := body.([]byte); !ok { // Check if body is of type []byte (Buffer)
			object["body"] = body
		} else {
			var jsonBody interface{}
			if err := json.Unmarshal(body.([]byte), &jsonBody); err == nil {
				object["body"] = jsonBody
			} else {
				object["body"] = string(body.([]byte)) // Convert []byte to string
			}
		}
	} else {
		object["body"] = "empty"
	}

	// Convert to JSON string and remove newlines
	jsonString, err := json.Marshal(object)
	if err != nil {
		logger.Errorf(ctx, "Error marshaling to JSON: %v\n", err)
		return ""
	}

	return string(jsonString)
}

func NewRestyClient(opts ...RestyOptions) *resty.Client {
	client := resty.New()

	// tracing instrumentation
	transport := http.DefaultTransport
	tracedTransport := otelhttp.NewTransport(transport)
	client.SetTransport(tracedTransport)

	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		ctx := r.Context()
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

		// Log the request
		// logger.Fields(map[string]interface{}{
		// 	logger.FIELD_OPERATOR_NAME: r.URL,
		// 	logger.FIELD_STEP_NAME:     "Request Method POST - WebClient",
		// }).Infof(ctx, MakeStringLogs(ctx, r.Header, r.Body, ""))

		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if r.StatusCode() == http.StatusOK {
			// ctx := r.Request.Context()

			// logger.Fields(map[string]interface{}{
			// 	logger.FIELD_OPERATOR_NAME: r.Request.URL,
			// 	logger.FIELD_STEP_NAME:     "Response Method POST - WebClient",
			// }).Infof(ctx, MakeStringLogs(ctx, r.Header(), r.Body(), r.Status()))

		} else {
			// ctx := r.Request.Context()
			// logger.Fields(map[string]interface{}{
			// 	logger.FIELD_OPERATOR_NAME: r.Request.URL,
			// 	logger.FIELD_STEP_NAME:     "Exception Method POST - WebClient",
			// }).Infof(ctx, string(r.Body()))
		}
		return nil
	})

	return client
}
