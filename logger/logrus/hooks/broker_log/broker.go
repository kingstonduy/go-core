package brokerLog

import (
	"context"
	"encoding/json"

	"github.com/kingstonduy/go-core/logger"
	logrusLib "github.com/kingstonduy/go-core/logger/logrus"
	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/transport/broker"
	"github.com/sirupsen/logrus"
)

var (
	DEFAULT_LOG_TOPIC = "mcs.logging.events"

	DEFAULT_PUBLISH_LOG_LEVELS = []logger.Level{
		logger.TraceLevel,
		logger.DebugLevel,
		logger.InfoLevel,
		logger.WarnLevel,
		logger.ErrorLevel,
		logger.FatalLevel,
	}

	DEFAULT_PUBLISH_EVENT_CONDITION BrokerLogEventCondition = func(event BrokerLogEvent) bool {
		published, ok := event[logger.FIELD_PUBLISHED].(bool)
		if !ok {
			return false
		}
		return published
	}
)

type BrokerLogEvent map[string]interface{}

type BrokerLogEventCondition = func(BrokerLogEvent) bool

type BrokerLogHookOptions struct {
	topic     string                  // log topic
	levels    []logger.Level          // event levels to be pushed
	condition BrokerLogEventCondition // fields conditions for an event to be pushed
}

type BrokerLogHookOption func(*BrokerLogHookOptions)

// define topic to push log event.
func WithBrokerLogHookTopic(topic string) BrokerLogHookOption {
	return func(o *BrokerLogHookOptions) {
		o.topic = topic
	}
}

// define levels to push log event.
func WithBrokerLogHookLevels(levels []logger.Level) BrokerLogHookOption {
	return func(o *BrokerLogHookOptions) {
		o.levels = levels
	}
}

// define condition of an log event's field. If matched, push it to broker topic.
func WithBrokerLogEventCondition(condition BrokerLogEventCondition) BrokerLogHookOption {
	return func(o *BrokerLogHookOptions) {
		o.condition = condition
	}
}

func newBrokerLogHookOptions(opts ...BrokerLogHookOption) BrokerLogHookOptions {
	options := BrokerLogHookOptions{
		topic:     DEFAULT_LOG_TOPIC,
		levels:    DEFAULT_PUBLISH_LOG_LEVELS,
		condition: DEFAULT_PUBLISH_EVENT_CONDITION,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

type BrokerLogHook struct {
	// Broker where log events published to
	Broker broker.Broker

	// Log Topic. Default: "mcs.logging.events"
	Topic string

	// Current service name
	ServiceName string

	// Log levels should be published. Default: all levels
	PublishLevels []logger.Level

	// Condition for log events, if matched, events will be published. Default: event["published"] == true
	PublishCondition BrokerLogEventCondition
}

func NewBrokerLogHook(logBroker broker.Broker, serviceName string, opts ...BrokerLogHookOption) logrus.Hook {
	options := newBrokerLogHookOptions(opts...)

	// connect if not already connected
	logBroker.Connect() //nolint

	return &BrokerLogHook{
		Broker:           logBroker,
		ServiceName:      serviceName,
		Topic:            options.topic,
		PublishLevels:    options.levels,
		PublishCondition: options.condition,
	}
}

// Fire implements logrus.Hook.
func (k *BrokerLogHook) Fire(entry *logrus.Entry) (err error) {
	// check event condition
	ctx, logEvent := k.getLogEvent(entry)

	// return if event is not matched the specified condition
	if !k.PublishCondition(*logEvent) {
		return nil
	}

	bytes, err := json.Marshal(logEvent)
	if err != nil {
		return err
	}
	// parse to broker message
	brokerMessage := broker.Message{
		Headers: map[string]string{},
		Body:    bytes,
	}

	return k.Broker.Publish(ctx, k.Topic, &brokerMessage)
}

// Levels implements logrus.Hook.
func (k *BrokerLogHook) Levels() []logrus.Level {
	levels := make([]logrus.Level, 0)

	for _, lv := range k.PublishLevels {
		levels = append(levels, logrusLib.LoggerToLogrusLevel(lv))
	}

	return levels
}

func (k *BrokerLogHook) getLogEvent(entry *logrus.Entry) (context.Context, *BrokerLogEvent) {
	ctx := entry.Context
	span := trace.ExtractSpanInfo(ctx)

	entryMessage, _ := entry.String() // ignore error when get string message

	logEvent := BrokerLogEvent{
		logger.FIELD_SERVICE:  k.ServiceName,
		logger.FIELD_TIME:     entry.Time,
		logger.FIELD_LEVEL:    logrusLib.LogrusToLoggerLevel(entry.Level).String(),
		logger.FIELD_MESSAGE:  entryMessage,
		logger.FIELD_TRACE_ID: span.TraceID,
		logger.FIELD_SPAN_ID:  span.SpanID,
	}

	for k, v := range entry.Data {
		logEvent[k] = v
	}

	return ctx, &logEvent
}
