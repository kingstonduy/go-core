package broker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kingstonduy/go-core/logger"
)

// Package broker is an interface used for asynchronous messaging

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Init(...BrokerOption) error
	Options() BrokerOptions
	Address() string
	Connect() error
	Disconnect() error
	Publish(ctx context.Context, topic string, m *Message, opts ...PublishOption) error
	PublishAndReceive(ctx context.Context, topic string, m *Message, opts ...PublishOption) (*Message, error)
	// With RabbitMQ: topic is a routing key. Using broker.WithSubscribeQueue(queue) to define the queue
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)

	String() string
}

// Handler is used to process messages via a subscription of a topic.
// The handler is passed a publication interface which contains the
// message and optional Ack method to acknowledge receipt of the message.
type Handler func(context.Context, Event) error

// Message is a message send/received from the broker.
type Message struct {
	Headers map[string]string
	Body    []byte
	Key     []byte
}

// Event is given to a subscription handler for processing.
type Event interface {
	// return event's topic
	Topic() string

	// return event's message
	Message() *Message

	// mark event as processed
	Ack() error

	// return error if event has error occurred
	Error() error

	// return timestamp of event
	Timestamp() time.Time
}

// Subscriber is a convenience return type for the Subscribe method.
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}

func MakeStringLogsKafka(ctx context.Context, bMsg Message) string {
	var parsedBody map[string]interface{}
	err := json.Unmarshal(bMsg.Body, &parsedBody)
	if err != nil {
		logger.Errorf(ctx, "Error unmarshaling to JSON: %v\n", err)
	}

	result := map[string]interface{}{
		"Headers": bMsg.Headers,
		"Key":     string(bMsg.Key), // Assuming the key is a string; adjust if needed
		"Body":    parsedBody,
	}

	// Convert the result to JSON
	finalMsg, err := json.Marshal(result)
	if err != nil {
		logger.Errorf(ctx, "Error marshaling final message: %v\n", err)
	}

	return string(finalMsg)
}
