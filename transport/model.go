package transport

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kingstonduy/go-core/errorx"
)

type Trace struct {
	From               string `json:"frm" xml:"frm"`
	To                 string `json:"to" xml:"to"`
	Cid                string `json:"cid" xml:"cid" validate:"required"`
	Sid                string `json:"sid" xml:"sid"`
	Cts                int64  `json:"cts" xml:"cts" validate:"required"`
	Sts                int64  `json:"sts" xml:"sts"`
	Dur                int64  `json:"dur" xml:"dur"`
	Username           string `json:"userName" xml:"userName"`
	ClientId           string `json:"clientId,omitempty" xml:"clientId,omitempty"`
	ReplyTo            string `json:"replyTo" xml:"replyTo"`
	TransactionTimeout int64  `json:"transactionTimeout,omitempty" xml:"transactionTimeout" validate:"max=120000"`
}

type Request[T any] struct {
	Trace Trace `json:"trace" xml:"trace" validate:"required"`
	Data  T     `json:"data" xml:"data"`
}

func (r *Request[T]) SetTrace(trace Trace) {
	r.Trace = trace
}

func (r *Request[T]) GetTrace() Trace {
	return r.Trace
}

type Result struct {
	StatusCode int         `json:"statusCode" xml:"statusCode"`
	Code       string      `json:"code" xml:"code"`
	Message    string      `json:"message" xml:"message"`
	Details    interface{} `json:"details" xml:"details" swaggertype:"object"`
}

type Response[T any] struct {
	Result Result `json:"result" xml:"result"`
	Trace  Trace  `json:"trace" xml:"trace"`
	Data   T      `json:"data" xml:"data"`
}

func (r *Response[T]) SetTrace(trace Trace) {
	r.Trace = trace
}

func (r *Response[T]) GetTrace() Trace {
	return r.Trace
}

var (
	DefaultSuccessResponse = Response[interface{}]{
		Result: Result{
			StatusCode: errorx.DefaultSuccessStatusCode,
			Code:       errorx.DefaultSuccessResponseCode,
			Message:    errorx.DefaultSuccessResponseMessage,
		},
		Trace: Trace{
			Sts: time.Now().UnixMilli(),
		},
		Data: nil,
	}

	DefaultFailureResponse = Response[interface{}]{
		Result: Result{
			StatusCode: errorx.DefaultFailureStatusCode,
			Code:       errorx.DefaultFailureResponseCode,
			Message:    errorx.DefaultFailureResponseMessage,
		},
		Trace: Trace{
			Sts: time.Now().UnixMilli(),
		},
		Data: nil,
	}
)

// Scan implements the sql.Scanner interface for Trace
func (t *Trace) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, t)
}

func (t *Result) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, t)
}

// Helper function to convert string or float64 to int64
func parseInt64Field(value interface{}, fieldName string) (int64, error) {
	if value == nil {
		return 0, nil
	}

	switch v := value.(type) {
	case string:
		intVal, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("error parsing %s field as string: %v", fieldName, err)
		}
		return intVal, nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type for %s field: %T", fieldName, v)
	}
}

// Custom UnmarshalJSON method to handle multiple fields as string or int64 in Trace
func (t *Trace) UnmarshalJSON(data []byte) error {
	// Create an alias to avoid infinite recursion
	type Alias Trace
	aux := &struct {
		Cts                interface{} `json:"cts"`
		Sts                interface{} `json:"sts"`
		Dur                interface{} `json:"dur"`
		TransactionTimeout interface{} `json:"transactionTimeout"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	// Unmarshal JSON into aux struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse each int64 field
	var err error
	if t.Cts, err = parseInt64Field(aux.Cts, "cts"); err != nil {
		return err
	}
	if t.Sts, err = parseInt64Field(aux.Sts, "sts"); err != nil {
		return err
	}
	if t.Dur, err = parseInt64Field(aux.Dur, "dur"); err != nil {
		return err
	}
	if t.TransactionTimeout, err = parseInt64Field(aux.TransactionTimeout, "transactionTimeout"); err != nil {
		return err
	}

	return nil
}

// Value implements the driver.Valuer interface for Trace
func (t Trace) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t Result) Value() (driver.Value, error) {
	return json.Marshal(t)
}

type WebClientResponseType struct {
	StatusCode    string
	StatusMessage string
	Headers       http.Header
	Body          []byte
}

// Option is a function that configures a WebClientResponseType.
type webClientOption func(*WebClientResponseType)

// WithStatusCode sets the status code.
func WithStatusCode(statusCode string) webClientOption {
	return func(w *WebClientResponseType) {
		w.StatusCode = statusCode
	}
}

// WithStatusMessage sets the status message.
func WithStatusMessage(statusMessage string) webClientOption {
	return func(w *WebClientResponseType) {
		w.StatusMessage = statusMessage
	}
}

// WithHeaders sets the headers.
func WithHeaders(headers http.Header) webClientOption {
	return func(w *WebClientResponseType) {
		w.Headers = headers
	}
}

// WithBody sets the response body.
func WithBody(body []byte) webClientOption {
	return func(w *WebClientResponseType) {
		w.Body = body
	}
}

// NewWebClientResponseType creates a new WebClientResponseType with the given options.
func NewWebClientResponseType(opts ...webClientOption) WebClientResponseType {
	response := WebClientResponseType{
		Headers: http.Header{}, // Initialize Headers to avoid nil
	}

	for _, opt := range opts {
		opt(&response)
	}

	return response
}
