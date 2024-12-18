package kafka

import (
	"context"

	"github.com/kingstonduy/go-core/transport/broker"

	"github.com/IBM/sarama"
)

var (
	DefaultProducerConfig = sarama.NewConfig()
	DefaultConsumerConfig = sarama.NewConfig()
)

type producerConfigKey struct{}
type consumerConfigKey struct{}

// Config for producers
func ProducerConfig(c *sarama.Config) broker.BrokerOption {
	return broker.SetBrokerOption(producerConfigKey{}, c)
}

// Config for producers
func ConsumerConfig(c *sarama.Config) broker.BrokerOption {
	return broker.SetBrokerOption(consumerConfigKey{}, c)
}

type subscribeContextKey struct{}

// SubscribeContext set the context for broker.SubscribeOption
func SubscribeContext(ctx context.Context) broker.SubscribeOption {
	return broker.SetSubscribeOption(subscribeContextKey{}, ctx)
}

type subscribeConfigKey struct{}

func SubscribeConfig(c *sarama.Config) broker.SubscribeOption {
	return broker.SetSubscribeOption(subscribeConfigKey{}, c)
}

type initialOffsetKey struct{}

// -1: OffsetNewest
// -2: OffsetOldest
// Default: OffsetNewest
func InitialOffset(offset int64) broker.SubscribeOption {
	return broker.SetSubscribeOption(initialOffsetKey{}, offset)
}

type asyncProduceErrorKey struct{}
type asyncProduceSuccessKey struct{}

func AsyncProducer(errors chan<- *sarama.ProducerError, successes chan<- *sarama.ProducerMessage) broker.BrokerOption {
	// set default opt
	var opt = func(options *broker.BrokerOptions) {}
	if successes != nil {
		opt = broker.SetBrokerOption(asyncProduceSuccessKey{}, successes)
	}
	if errors != nil {
		opt = broker.SetBrokerOption(asyncProduceErrorKey{}, errors)
	}
	return opt
}
