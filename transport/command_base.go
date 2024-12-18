package transport

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kingstonduy/go-core/validation"
)

// @Author: Nghiant5
// This Command Base is used for outbox pattern entity
type Command struct {
	CommandID            string  `json:"COMMAND_ID"`
	AggregateID          string  `json:"AGGREGATE_ID"`
	CommandType          string  `json:"COMMAND_TYPE"`
	AggregateType        string  `json:"AGGREGATE_TYPE"`
	Payload              string  `json:"PAYLOAD"`
	ProcessedAt          int64   `json:"PROCESSED_AT"`
	ReplyTo              string  `json:"REPLY_TO"`
	TransactionCreatedAt int64   `json:"TRANSACTION_CREATED_AT"`
	Trace                Trace   `json:"TRACE"`
	Result               *Result `json:"RESULT"`
	Forward              string  `json:"FORWARD"`
}

func NewDefaultCommand() Command {
	return Command{
		Trace: Trace{
			TransactionTimeout: 60000,
		},
	}
}

// NewBaseCommand new base Command constructor with configured CommandID, Aggregate properties and Timestamp.
func NewBaseCommand(aggregateId string, aggregateType string, CommandType string,
	replyTo string, transactionCreatedAt int64, trace Trace, result *Result, forward string) *Command {
	return &Command{
		CommandID:            uuid.New().String(),
		AggregateType:        aggregateType,
		AggregateID:          aggregateId,
		CommandType:          CommandType,
		ProcessedAt:          time.Now().UnixMilli(),
		ReplyTo:              replyTo,
		TransactionCreatedAt: transactionCreatedAt,
		Trace:                trace,
		Result:               result,
		Forward:              forward,
	}
}

func NewCommand(aggregateId string, aggregateType string, CommandType string, data string,
	replyTo string, transactionCreatedAt int64, trace Trace, result *Result, forward string) *Command {
	return &Command{
		CommandID:            uuid.New().String(),
		AggregateID:          aggregateId,
		CommandType:          CommandType,
		AggregateType:        aggregateType,
		Payload:              data,
		ProcessedAt:          time.Now().UnixMilli(),
		ReplyTo:              replyTo,
		TransactionCreatedAt: transactionCreatedAt,
		Trace:                trace,
		Result:               result,
		Forward:              forward,
	}
}

func (e *Command) SetNewCommandId() {
	e.CommandID = uuid.New().String()
}

// GetCommandID get CommandID of the Command.
func (e *Command) GetCommandID() string {
	return e.CommandID
}

// GetData The data attached to the Command serialized to bytes.
func (e *Command) GetPayload() interface{} {
	return e.Payload
}

// SetData add the data attached to the Command serialized to bytes.
func (e *Command) SetPayload(data string) {
	e.Payload = data

}

func (e *Command) GetJsonPayload(data interface{}) error {
	err := json.Unmarshal([]byte(e.Payload), data)
	if err != nil {
		return err
	}
	return validation.Validate(data)
}

// GetCommandType returns the CommandType of the Command.
func (e *Command) GetCommandType() string {
	return e.CommandType
}

// GetAggregateType is the AggregateType that the Command can be applied to.
func (e *Command) GetAggregateType() string {
	return e.AggregateType
}

// SetAggregateType set the AggregateType that the Command can be applied to.
func (e *Command) SetAggregateType(aggregateType string) {
	e.AggregateType = aggregateType
}

// GetAggregateID is the AggregateID of the Aggregate that the Command belongs to
func (e *Command) GetAggregateID() string {
	return e.AggregateID
}

// SetAggregateType set the ReplyTo that the Command can be applied to.
func (e *Command) SetReplyTo(replyTo string) {
	e.ReplyTo = replyTo
}

// GetAggregateID is the AggregateID of the Aggregate that the Command belongs to
func (e *Command) GetReplyTo() string {
	return e.ReplyTo
}

// GetString A string representation of the Command.
func (e *Command) GetString() string {
	return fmt.Sprintf("Command: %+v", e)
}

func (e *Command) SetTransactionCreatedAt(time int64) {
	e.TransactionCreatedAt = time
}

func (e *Command) GetTransactionCreatedAt() int64 {
	return e.TransactionCreatedAt
}

func (e *Command) String() string {
	return fmt.Sprintf("(Command) AggregateID: %v, CommandType: %v, AggregateType: %v, ProcessAt: %v, CommandID: %v, Payload: %v, Trace: %+v, Result: %+v",
		e.AggregateID,
		e.CommandType,
		e.AggregateType,
		e.ProcessedAt,
		e.CommandID,
		e.Payload,
		e.Trace,
		e.Result,
	)
}

func (e *Command) StringNoPayload() string {
	return fmt.Sprintf("(Command) AggregateID: %v, CommandType: %v, AggregateType: %v, ProcessAt: %v, CommandID: %v, Trace: %+v, Result: %+v",
		e.AggregateID,
		e.CommandType,
		e.AggregateType,
		e.ProcessedAt,
		e.CommandID,
		e.Trace,
		e.Result,
	)
}
