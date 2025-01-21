package cmd_pipeline

import (
	"encoding/json"

	"github.com/kingstonduy/go-core/transport"
)

type Outbox struct {
	AggregateID string `json:"AGGREGATE_ID"`
	CommandID   string `json:"COMMAND_ID"`
	CommandType string `json:"COMMAND_TYPE"`
	Payload     string `json:"PAYLOAD"`
	Trace       string `json:"TRACE"`
	ReplyTo     string `json:"REPLY_TO" db:"REPLY_TO"`
	TraceParent string `json:"TRACE_PARENT" db:"TRACE_PARENT"`
}

type OutboxWithTrace struct {
	AggregateID string          `json:"AGGREGATE_ID"`
	CommandID   string          `json:"COMMAND_ID"`
	CommandType string          `json:"COMMAND_TYPE"`
	Payload     string          `json:"PAYLOAD"`
	Trace       transport.Trace `json:"TRACE"`
	ReplyTo     string          `json:"REPLY_TO" db:"REPLY_TO"`
	TraceParent string          `json:"TRACE_PARENT" db:"TRACE_PARENT"`
}

func (o *OutboxWithTrace) ToString() string {
	s, _ := json.Marshal(o)
	return string(s)
}

func (o *Outbox) ToString() string {
	s, _ := json.Marshal(o)
	return string(s)
}

func (o *Outbox) ToOutboxWithTrace() OutboxWithTrace {
	res := OutboxWithTrace{
		AggregateID: o.AggregateID,
		CommandID:   o.CommandID,
		CommandType: o.CommandType,
		Payload:     o.Payload,
		Trace:       transport.Trace{},
		ReplyTo:     o.ReplyTo,
		TraceParent: o.TraceParent,
	}

	json.Unmarshal([]byte(o.Trace), &res.Trace)
	return res
}
