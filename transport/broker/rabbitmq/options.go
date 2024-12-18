package rabbitmq

import (
	"context"
	"time"

	"github.com/kingstonduy/go-core/transport/broker"
)

type durableQueueKey struct{}
type headersKey struct{}
type queueArgumentsKey struct{}
type prefetchCountKey struct{}
type prefetchGlobalKey struct{}
type confirmPublishKey struct{}
type exchangeKey struct{}
type exchangeTypeKey struct{}
type withoutExchangeKey struct{}
type requeueOnErrorKey struct{}
type deliveryMode struct{}
type priorityKey struct{}
type contentType struct{}
type contentEncoding struct{}
type correlationID struct{}
type replyTo struct{}
type expiration struct{}
type messageID struct{}
type timestamp struct{}
type typeMsg struct{}
type userID struct{}
type appID struct{}
type externalAuth struct{}
type durableExchange struct{}

/*
	DefaultWithoutExchange = false
	// The amqp library does not seem to set these when using amqp.DialConfig
	// (even though it says so in the comments) so we set them manually to make
	// sure to not brake any existing functionality.
	defaultHeartbeat = 10 * time.Second
	defaultLocale    = "en_US"
*/

// ================= BROKER OPTIONS =================
// DurableExchange is an option to set the Exchange to be durable.
// Default false
func DurableExchange() broker.BrokerOption {
	return broker.SetBrokerOption(durableExchange{}, true)
}

// ExchangeName is an option to set the ExchangeName.
func ExchangeName(e string) broker.BrokerOption {
	return broker.SetBrokerOption(exchangeKey{}, e)
}

func ExternalAuth() broker.BrokerOption {
	return broker.SetBrokerOption(externalAuth{}, ExternalAuthentication{})
}

// ExchangeType is an option to set the rabbitmq exchange type.
func ExchangeType(t MQExchangeType) broker.BrokerOption {
	return broker.SetBrokerOption(exchangeTypeKey{}, t)
}

// WithoutExchange is an option to use the rabbitmq default exchange.
// means it would not create any custom exchange.
// DefaultWithoutExchange = false
func WithoutExchange() broker.BrokerOption {
	return broker.SetBrokerOption(withoutExchangeKey{}, true)
}

// PrefetchCount ...
// DefaultPrefetchCount   = 0
func PrefetchCount(c int) broker.BrokerOption {
	return broker.SetBrokerOption(prefetchCountKey{}, c)
}

// PrefetchGlobal creates a durable queue when subscribing.
// DefaultPrefetchGlobal  = false
func PrefetchGlobal() broker.BrokerOption {
	return broker.SetBrokerOption(prefetchGlobalKey{}, true)
}

// ConfirmPublish ensures all published messages are confirmed by waiting for an ack/nack from the broker.
// DefaultConfirmPublish  = false
func ConfirmPublish() broker.BrokerOption {
	return broker.SetBrokerOption(confirmPublishKey{}, true)
}

// ================= SUBSCRIBE OPTIONS =================
// DurableQueue creates a durable queue when subscribing.
func DurableQueue() broker.SubscribeOption {
	return broker.SetSubscribeOption(durableQueueKey{}, true)
}

// Headers adds headers used by the headers exchange.
func Headers(h map[string]interface{}) broker.SubscribeOption {
	return broker.SetSubscribeOption(headersKey{}, h)
}

// QueueArguments sets arguments for queue creation.
func QueueArguments(h map[string]interface{}) broker.SubscribeOption {
	return broker.SetSubscribeOption(queueArgumentsKey{}, h)
}

// DefaultRequeueOnError  = false
func RequeueOnError() broker.SubscribeOption {
	return broker.SetSubscribeOption(requeueOnErrorKey{}, true)
}

type subscribeContextKey struct{}

// SubscribeContext set the context for broker.SubscribeOption.
func SubscribeContext(ctx context.Context) broker.SubscribeOption {
	return broker.SetSubscribeOption(subscribeContextKey{}, ctx)
}

type ackSuccessKey struct{}

// AckOnSuccess will automatically acknowledge messages when no error is returned.
func AckOnSuccess() broker.SubscribeOption {
	return broker.SetSubscribeOption(ackSuccessKey{}, true)
}

// ================= PUBLISH OPTIONS =================

// DeliveryMode sets a delivery mode for publishing.
// PublishDeliveryMode client.PublishOption for setting message "delivery mode"
// mode , Transient (0 or 1) or Persistent (2)
func DeliveryMode(value uint8) broker.PublishOption {
	return broker.SetPublishOption(deliveryMode{}, value)
}

// Priority sets a priority level for publishing.
func Priority(value uint8) broker.PublishOption {
	return broker.SetPublishOption(priorityKey{}, value)
}

// ContentType sets a property MIME content type for publishing.
func ContentType(value string) broker.PublishOption {
	return broker.SetPublishOption(contentType{}, value)
}

// ContentEncoding sets a property MIME content encoding for publishing.
func ContentEncoding(value string) broker.PublishOption {
	return broker.SetPublishOption(contentEncoding{}, value)
}

// CorrelationID sets a property correlation ID for publishing.
func CorrelationID(value string) broker.PublishOption {
	return broker.SetPublishOption(correlationID{}, value)
}

// ReplyTo sets a property address to to reply to (ex: RPC) for publishing.
func ReplyTo(value string) broker.PublishOption {
	return broker.SetPublishOption(replyTo{}, value)
}

// Expiration sets a property message expiration spec for publishing.
func Expiration(value string) broker.PublishOption {
	return broker.SetPublishOption(expiration{}, value)
}

// MessageId sets a property message identifier for publishing.
func MessageId(value string) broker.PublishOption {
	return broker.SetPublishOption(messageID{}, value)
}

// Timestamp sets a property message timestamp for publishing.
func Timestamp(value time.Time) broker.PublishOption {
	return broker.SetPublishOption(timestamp{}, value)
}

// TypeMsg sets a property message type name for publishing.
func TypeMsg(value string) broker.PublishOption {
	return broker.SetPublishOption(typeMsg{}, value)
}

// UserID sets a property user id for publishing.
func UserID(value string) broker.PublishOption {
	return broker.SetPublishOption(userID{}, value)
}

// AppID sets a property application id for publishing.
func AppID(value string) broker.PublishOption {
	return broker.SetPublishOption(appID{}, value)
}
