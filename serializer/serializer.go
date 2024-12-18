package serializer

import (
	"encoding/json"

	"github.com/kingstonduy/go-core/transport"
	"github.com/kingstonduy/go-core/validation"
	"github.com/pkg/errors"
)

type SerialOptions struct {
	result  *transport.Result
	payload interface{}
}

type SerialOption func(s *SerialOptions)

// The request header mapped to response headers
func WithResult(result *transport.Result) SerialOption {
	return func(options *SerialOptions) {
		options.result = result
	}
}

// The request header mapped to response headers
func WithPayload(payload interface{}) SerialOption {
	return func(options *SerialOptions) {
		options.payload = payload
	}
}

// @Author: Nghiant5
// This is used for Serilize and Deserialize Command
var (
	ErrInvalidcommand = errors.New("invalid command")
)

func SerializeCommand(aggregateId string, aggregateType string, commandType string,
	transactionCreatedAt int64, trace transport.Trace, replyTo string, forward string, opts ...SerialOption) (*transport.Command, error) {
	options := SerialOptions{}

	for _, opt := range opts {
		opt(&options)
	}
	var payload string
	if options.payload != nil {
		commandBytes, err := json.Marshal(options.payload)
		if err != nil {
			return nil, errors.Wrapf(err, "serializer.Marshal aggregateID: %s", aggregateId)
		}
		payload = string(commandBytes)
	}

	if options.result != nil {
		return transport.NewCommand(aggregateId, aggregateType, commandType, payload, replyTo, transactionCreatedAt, trace, options.result, forward), nil
	}
	return transport.NewCommand(aggregateId, aggregateType, commandType, payload, replyTo, transactionCreatedAt, trace, nil, forward), nil

}

// Deserialize will response a data type of payload that use passed to
// And it automatically validate data with the field which is tagged
func DeserializeCommand[T any](command transport.Command) (*T, error) {
	var commandTarget T
	if err := validation.Validate(command.Trace); err != nil {
		return nil, err
	}
	return deserializeCommand[T](command, commandTarget)
}

func deserializeCommand[T any](command transport.Command, targetcommand T) (*T, error) {
	if err := command.GetJsonPayload(&targetcommand); err != nil {
		return nil, errors.Wrapf(err, "command.GetJsonData type: %s", command.CommandType)
	}
	return &targetcommand, nil
}
