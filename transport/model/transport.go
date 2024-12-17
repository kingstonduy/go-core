package transport

import (
	"time"

	errorx "github.com/kingstonduy/go-core/error"
)

// Event represents a generic event with a payload of type T.
type Event[T any] struct {
	AggregateID   string `json:"aggregateID,omitempty"`
	AggregateType string `json:"aggregateType,omitempty"`
	EventID       string `json:"eventID,omitempty"`
	EventType     string `json:"eventType,omitempty"`
	PayLoad       T      `json:"payLoad,omitempty"`
	Trace         Trace  `json:"trace,omitempty"`
	Result        Result `json:"result,omitempty"`
}

type RequestType[T any] struct {
	Trace  Trace  `json:"trace,omitempty" validate:"required"`
	Data   T      `json:"data,omitempty"`
	Result Result `json:"result,omitempty"`
}

type ResponseType[T any] struct {
	Result Result `json:"result,omitempty"`
	Trace  Trace  `json:"trace,omitempty" validate:"required"`
	Data   T      `json:"data,omitempty"`
}

// Trace represents tracing information for the event.
type Trace struct {
	Cid     string `json:"cid,omitempty"`
	Cts     int64  `json:"cts,omitempty"`
	Sid     string `json:"sid,omitempty"`
	Sts     int64  `json:"sts,omitempty"`
	From    string `json:"frm,omitempty"`
	To      string `json:"to,omitempty"`
	Dur     int64  `json:"dur,omitempty"`
	TimeOut int64  `json:"timeout,omitempty"`
}

// Result represents the result of processing an event.
type Result struct {
	Code    string `json:"code,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func GetResultFromErrorx(errx errorx.Errorx) Result {
	return Result{Code: errx.Code, Status: errx.Status, Message: errx.Message, Detail: errx.Detail}
}

// return http response trace from request trace
func GetResponseTrace(reqTrace Trace) Trace {
	now := time.Now()
	return Trace{
		Cid:     reqTrace.Cid,
		Cts:     reqTrace.Cts,
		Sid:     reqTrace.Sid,
		Sts:     now.Unix(),
		From:    reqTrace.To,
		To:      reqTrace.From,
		TimeOut: reqTrace.TimeOut,
		Dur:     now.Unix() - reqTrace.Cts,
	}
}
