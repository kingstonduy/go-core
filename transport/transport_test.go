package transport

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/logger/logrus"
	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/trace/otel"
	"github.com/kingstonduy/go-core/validation"
	"github.com/kingstonduy/go-core/validation/goplayaround"
	"github.com/stretchr/testify/assert"
)

type DataRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type DataResponse struct {
	Data string `json:"data"`
}

var (
	request = Request[DataRequest]{
		Trace: Trace{
			Cid:      "123123123",
			Sid:      "1234133423",
			From:     "10.0.0.2",
			To:       "1.2.3.4",
			Cts:      1711592361558,
			Username: "username",
			// MessageType: "CreateUserRequest",
			// ReplyTo:     []string{"topic.1", "topic.2"},
		},
		Data: DataRequest{
			Name: "OCB Staff",
			Age:  22,
		},
	}
)

func getValidator(t *testing.T) validation.Validator {
	return goplayaround.NewGpValidator()
}

func getTracer(t *testing.T) trace.Tracer {
	tracer, err := otel.NewOpenTelemetryTracer(
		context.Background(),
		trace.WithTraceExporterEndpoint("localhost:4318"),
	)
	if err != nil {
		t.Error(err)
		return nil
	}

	return tracer
}

func getLogger(t *testing.T, tracer trace.Tracer) logger.Logger {
	return logrus.NewLogrusLogger()
}

func TestTransport(t *testing.T) {
	t.Run("Test monitoring request", func(t *testing.T) {
		TestMonitorRequest(t)
		TestMonitorRequestWithTracingLogging(t)
		TestMonitorRequestWithSpanInfo(t)
	})

	t.Run("Test responses", func(t *testing.T) {
		TestGetResponseSuccess(t)
		TestGetResponseError(t)
		TestGetResponseSuccessEmpty(t)
		TestGetResponseErrorsEmpty(t)
	})

	t.Run("Test with default tracer and logger", func(t *testing.T) {
		TestGetResponseDefaultError(t)
		TestGetResponseWithDefault(t)
		TestGetResponseCustomErrorWithDefault(t)
		TestGetResponseErrorWithDefault(t)
	})

	t.Run("Test message type and reply to", func(t *testing.T) {
		TestMessageTypeAndReplyTo(t)
	})
}

func TestMonitorRequest(t *testing.T) {
	validator := goplayaround.NewGpValidator()
	err := validator.Validate(request)
	if err != nil {
		t.Error(err)
	}

	ctx := MonitorRequest(
		context.Background(),
		MonitorRequestData{
			Request: request,
		},
	)

	time.Sleep(200 * time.Millisecond)
	assert.NotNil(t, ctx)
}

func TestMonitorRequestWithTracingLogging(t *testing.T) {
	validator := getValidator(t)
	err := validator.Validate(request)
	if err != nil {
		t.Error(err)
	}

	tracer := getTracer(t)
	logger := getLogger(t, tracer)

	ctx := MonitorRequest(
		context.Background(),
		MonitorRequestData{
			Request: request,
		},
		WithTracer(tracer),
		WithLogger(logger),
	)

	time.Sleep(200 * time.Millisecond)

	assert.NotNil(t, ctx)
	assert.NotNil(t, ctx.Value(trace.SpanInfoKey{}))
}

func TestMonitorRequestWithSpanInfo(t *testing.T) {
	validator := getValidator(t)
	err := validator.Validate(request)
	if err != nil {
		t.Error(err)
	}

	tracer := getTracer(t)
	logger := getLogger(t, tracer)
	data := MonitorRequestData{
		From:          "ABC",
		To:            "ABC",
		ClientIP:      "456.456.456.456",
		Protocol:      "HTTP",
		Method:        "POST",
		ServiceDomain: "servicename",
		RequestPath:   "/api/v1/get",
		UserAgent:     "postman",
		ClientTime:    time.Now().UnixMilli(),
		Hostname:      "pod-1",
		TransactionID: "123456123456",
		ContentLength: 120,
		Request:       request,
		ClientID:      uuid.New().String(),
	}

	ctx := MonitorRequest(
		context.Background(),
		data,
		WithTracer(tracer),
		WithLogger(logger),
	)

	time.Sleep(200 * time.Millisecond)

	span := tracer.ExtractSpanInfo(ctx)

	assert.NotNil(t, span.TraceID)
	assert.NotNil(t, span.SpanID)
	assert.Equal(t, data.ClientIP, span.ClientIP)
	assert.Equal(t, data.Protocol, span.Protocol)
	assert.Equal(t, data.Method, span.Method)
	assert.Equal(t, data.ServiceDomain, span.ServiceDomain)
	assert.Equal(t, data.RequestPath, span.RequestPath)
	assert.Equal(t, data.UserAgent, span.UserAgent)
	assert.Equal(t, data.ClientTime, span.ClientTime)
	assert.Equal(t, data.Hostname, span.Hostname)
	assert.Equal(t, data.ClientID, span.ClientID)
	assert.Equal(t, data.From, span.From)
	assert.Equal(t, data.To, span.To)
	assert.Equal(t, data.TransactionID, span.TransactionID)
}

func TestGetResponseSuccess(t *testing.T) {
	data := MonitorRequestData{
		From:          "ABC",
		To:            "ABC",
		ClientIP:      "456.456.456.456",
		Protocol:      "HTTP",
		Method:        "POST",
		ServiceDomain: "servicename",
		RequestPath:   "/api/v1/get",
		UserAgent:     "postman",
		ClientTime:    time.Now().UnixMilli(),
		Hostname:      "pod-1",
		TransactionID: "123456123456",
		ContentLength: 120,
		Request:       request,
		ClientID:      uuid.New().String(),
		Username:      "username",
	}

	tracer := getTracer(t)
	logger := getLogger(t, tracer)

	ctx := MonitorRequest(
		context.Background(),
		data,
		WithTracer(tracer),
		WithLogger(logger),
	)

	time.Sleep(200 * time.Millisecond)

	res := GetResponse[*DataResponse](
		ctx,
		WithData(&DataResponse{Data: "Response data"}),
		WithTraceExtractor(tracer),
	)

	j, _ := json.Marshal(res)
	t.Logf("Response: %v", string(j))

	assert.Equal(t, 200, res.Result.StatusCode)
	assert.Equal(t, "00", res.Result.Code)
	assert.Equal(t, "Successful", res.Result.Message)
	assert.Greater(t, res.Trace.Dur, int64(0))
	assert.Equal(t, res.Trace.Cid, data.ClientID)
	assert.Equal(t, res.Trace.Cts, data.ClientTime)
	assert.Equal(t, res.Trace.From, data.From)
	assert.Equal(t, res.Trace.To, data.To)
	assert.Equal(t, res.Trace.Username, data.Username)
	assert.NotNil(t, res.Trace.Sid)
}

func TestGetResponseError(t *testing.T) {
	data := MonitorRequestData{
		From:          "ABC",
		To:            "ABC",
		ClientIP:      "456.456.456.456",
		Protocol:      "HTTP",
		Method:        "POST",
		ServiceDomain: "servicename",
		RequestPath:   "/api/v1/get",
		UserAgent:     "postman",
		ClientTime:    time.Now().UnixMilli(),
		Hostname:      "pod-1",
		TransactionID: "123456123456",
		ContentLength: 120,
		Request:       request,
		ClientID:      uuid.New().String(),
	}

	err := errorx.NewError(99, "99", "error: %s", "test")
	tracer := getTracer(t)
	logger := getLogger(t, tracer)

	ctx := MonitorRequest(
		context.Background(),
		data,
		WithTracer(tracer),
		WithLogger(logger),
	)

	time.Sleep(200 * time.Millisecond)

	res := GetResponse[*DataResponse](
		ctx,
		WithData(&DataResponse{Data: "Response data"}),
		WithTraceExtractor(tracer),
		WithTrace(&request.Trace),
		WithError(err),
	)

	j, _ := json.Marshal(res)
	t.Logf("Response: %v", string(j))

	assert.Equal(t, 99, res.Result.StatusCode)
	assert.Equal(t, "99", res.Result.Code)
	assert.Equal(t, "error: test", res.Result.Message)
	assert.Greater(t, res.Trace.Dur, int64(0))
	assert.Equal(t, res.Trace.Cid, request.Trace.Cid)
	assert.Equal(t, res.Trace.Sid, request.Trace.Sid)
	assert.Equal(t, res.Trace.Cts, request.Trace.Cts)
	assert.Equal(t, res.Trace.From, request.Trace.From)
	assert.Equal(t, res.Trace.From, request.Trace.From)
	assert.Equal(t, res.Trace.To, request.Trace.To)
	assert.Equal(t, res.Trace.Username, request.Trace.Username)
}

func TestGetResponseDefaultError(t *testing.T) {
	data := MonitorRequestData{
		From:          "ABC",
		To:            "ABC",
		ClientIP:      "456.456.456.456",
		Protocol:      "HTTP",
		Method:        "POST",
		ServiceDomain: "servicename",
		RequestPath:   "/api/v1/get",
		UserAgent:     "postman",
		ClientTime:    time.Now().UnixMilli(),
		Hostname:      "pod-1",
		TransactionID: "123456123456",
		ContentLength: 120,
		Request:       request,
		ClientID:      uuid.New().String(),
	}

	err := errors.New("error: test")
	tracer := getTracer(t)
	logger := getLogger(t, tracer)

	ctx := MonitorRequest(
		context.Background(),
		data,
		WithTracer(tracer),
		WithLogger(logger),
	)

	time.Sleep(200 * time.Millisecond)

	res := GetResponse[*DataResponse](
		ctx,
		WithTrace(&request.Trace),
		WithError(err),
	)

	j, _ := json.Marshal(res)
	t.Logf("Response: %v", string(j))

	assert.Equal(t, errorx.DefaultFailureStatusCode, res.Result.StatusCode)
	assert.Equal(t, errorx.DefaultFailureResponseCode, res.Result.Code)
	assert.Contains(t, res.Result.Message, err.Error())
	assert.Greater(t, res.Trace.Dur, int64(0))
	assert.Equal(t, res.Trace.Cid, request.Trace.Cid)
	assert.Equal(t, res.Trace.Sid, request.Trace.Sid)
	assert.Equal(t, res.Trace.Cts, request.Trace.Cts)
	assert.Equal(t, res.Trace.From, request.Trace.From)
	assert.Equal(t, res.Trace.From, request.Trace.From)
	assert.Equal(t, res.Trace.To, request.Trace.To)
	assert.Equal(t, res.Trace.Username, request.Trace.Username)
}

func TestGetResponseSuccessEmpty(t *testing.T) {
	res := GetResponse[*DataResponse](context.Background())

	j, _ := json.Marshal(res)
	t.Logf("Response: %v", string(j))

	assert.Equal(t, errorx.DefaultSuccessStatusCode, res.Result.StatusCode)
	assert.Equal(t, errorx.DefaultSuccessResponseCode, res.Result.Code)
	assert.Equal(t, errorx.DefaultSuccessResponseMessage, res.Result.Message)
}

func TestGetResponseErrorsEmpty(t *testing.T) {
	res := GetResponse[*DataResponse](
		context.Background(),
		WithError(errors.New("error: test")),
	)

	j, _ := json.Marshal(res)
	t.Logf("Response: %v", string(j))

	assert.Equal(t, errorx.DefaultFailureStatusCode, res.Result.StatusCode)
	assert.Equal(t, errorx.DefaultFailureResponseCode, res.Result.Code)
	assert.Contains(t, res.Result.Message, "error: test")
}

func TestGetResponseWithDefault(t *testing.T) {
	data := MonitorRequestData{
		From:          "ABC",
		To:            "ABC",
		ClientIP:      "456.456.456.456",
		Protocol:      "HTTP",
		Method:        "POST",
		ServiceDomain: "servicename",
		RequestPath:   "/api/v1/get",
		UserAgent:     "postman",
		ClientTime:    time.Now().UnixMilli(),
		Hostname:      "pod-1",
		TransactionID: "123456123456",
		ContentLength: 120,
		Request:       request,
		ClientID:      uuid.New().String(),
		Username:      "username",
	}

	trace.SetDefaultTracer(getTracer(t))
	logger.SetDefaultLogger(getLogger(t, nil))

	dataRes := DataResponse{Data: "Response data"}
	ctx := context.Background()
	ctx = MonitorRequest(ctx, data)

	// response duration
	time.Sleep(1000 * time.Millisecond)

	res := GetResponse[*DataResponse](
		ctx,
		WithData(&dataRes),
		WithError(errors.New("error: test")),
	)

	j, _ := json.Marshal(res)
	t.Logf("Response: %v", string(j))

	assert.Equal(t, errorx.DefaultFailureStatusCode, res.Result.StatusCode)
	assert.Equal(t, errorx.DefaultFailureResponseCode, res.Result.Code)
	assert.Contains(t, res.Result.Message, "error: test")
	assert.Equal(t, dataRes.Data, res.Data.Data)

	time.Sleep(100 * time.Millisecond)
}

func TestGetResponseErrorWithDefault(t *testing.T) {
	data := MonitorRequestData{
		From:          "ABC",
		To:            "ABC",
		ClientIP:      "456.456.456.456",
		Protocol:      "HTTP",
		Method:        "POST",
		ServiceDomain: "servicename",
		RequestPath:   "/api/v1/get",
		UserAgent:     "postman",
		ClientTime:    time.Now().UnixMilli(),
		Hostname:      "pod-1",
		TransactionID: "123456123456",
		ContentLength: 120,
		Request:       request,
		ClientID:      uuid.New().String(),
		Username:      "username",
		XForwardedFor: "x-forwarded-for",
		RemoteHost:    "remote-host",
	}

	trace.SetDefaultTracer(getTracer(t))
	logger.SetDefaultLogger(getLogger(t, nil))

	ctx := context.Background()
	ctx = MonitorRequest(ctx, data)
	res := GetResponse[*DataResponse](
		ctx,
		WithData(&DataResponse{Data: "Response data"}),
	)

	j, _ := json.Marshal(res)
	t.Logf("Response: %v", string(j))

	assert.Equal(t, errorx.DefaultSuccessStatusCode, res.Result.StatusCode)
	assert.Equal(t, errorx.DefaultSuccessResponseCode, res.Result.Code)
	assert.Equal(t, errorx.DefaultSuccessResponseMessage, res.Result.Message)
}

func TestGetResponseCustomErrorWithDefault(t *testing.T) {
	data := MonitorRequestData{
		From:          "ABC",
		To:            "ABC",
		ClientIP:      "456.456.456.456",
		Protocol:      "HTTP",
		Method:        "POST",
		ServiceDomain: "servicename",
		RequestPath:   "/api/v1/get",
		UserAgent:     "postman",
		ClientTime:    time.Now().UnixMilli(),
		Hostname:      "pod-1",
		TransactionID: "123456123456",
		ContentLength: 120,
		Request:       request,
		ClientID:      uuid.New().String(),
		Username:      "username",
	}

	trace.SetDefaultTracer(getTracer(t))
	logger.SetDefaultLogger(getLogger(t, nil))

	// response duration
	time.Sleep(100 * time.Millisecond)

	ctx := context.Background()
	ctx = MonitorRequest(ctx, data)
	res := GetResponse[*DataResponse](
		ctx,
		WithData(&DataResponse{Data: "Response data"}),
		WithError(errorx.NewError(99, "99", "error: %s", "test")),
	)

	j, _ := json.Marshal(res)
	t.Logf("Response: %v", string(j))
	time.Sleep(100 * time.Millisecond)
}

func TestMessageTypeAndReplyTo(t *testing.T) {
	data := MonitorRequestData{
		From:          "ABC",
		To:            "ABC",
		ClientIP:      "456.456.456.456",
		Protocol:      "HTTP",
		Method:        "POST",
		ServiceDomain: "servicename",
		RequestPath:   "/api/v1/get",
		UserAgent:     "postman",
		ClientTime:    time.Now().UnixMilli(),
		Hostname:      "pod-1",
		TransactionID: "123456123456",
		ContentLength: 120,
		Request:       request,
		ClientID:      uuid.New().String(),
		Username:      "username",
		MessageType:   "CreateUserRequest",
		ReplyTo:       []string{"queue.1", "queue.2"},
	}

	tracer := getTracer(t)
	trace.SetDefaultTracer(tracer)
	logger.SetDefaultLogger(getLogger(t, tracer))

	ctx := MonitorRequest(
		context.Background(),
		data,
	)

	time.Sleep(200 * time.Millisecond)

	res := GetResponse[*DataResponse](
		ctx,
		WithData(&DataResponse{Data: "Response data"}),
		WithTraceExtractor(tracer),
	)

	j, _ := json.Marshal(res)
	t.Logf("Response: %v", string(j))

	assert.Equal(t, 200, res.Result.StatusCode)
	assert.Equal(t, "00", res.Result.Code)
	assert.Equal(t, "Successful", res.Result.Message)
	assert.Greater(t, res.Trace.Dur, int64(0))
	assert.Equal(t, res.Trace.Cid, data.ClientID)
	assert.Equal(t, res.Trace.Cts, data.ClientTime)
	assert.Equal(t, res.Trace.From, data.From)
	assert.Equal(t, res.Trace.To, data.To)
	assert.Equal(t, res.Trace.Username, data.Username)
	assert.NotNil(t, res.Trace.Sid)
	// assert.NotNil(t, res.Trace.MessageType)
	// assert.NotNil(t, res.Trace.MessageType, data.MessageType)
	// assert.NotNil(t, res.Trace.ReplyTo, data.ReplyTo)
}

func TestNewClientId(t *testing.T) {
	trace := Trace{
		ClientId: "data",
	}

	assert.Equal(t, trace.ClientId, "data")
}
