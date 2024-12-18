package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/transport/broker"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"go.opentelemetry.io/otel"
)

var (
	RequestReplyTimeout = time.Second * 60
)

type kBroker struct {
	addrs []string

	client         sarama.Client        // broker connection client
	syncProducer   sarama.SyncProducer  // sync producer
	asyncProducer  sarama.AsyncProducer // async producer
	consumerGroups []sarama.ConsumerGroup
	connected      bool
	scMutex        sync.Mutex
	opts           broker.BrokerOptions

	// request-reply patterns
	resps           sync.Map
	respSubscribers sync.Map

	codec Codec
}

func NewKafkaBroker(opts ...broker.BrokerOption) broker.Broker {
	options := broker.NewBrokerOptions(opts...)

	var cAddrs []string
	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}

	if len(cAddrs) == 0 {
		cAddrs = []string{DefaultKafkaBroker}
	}

	return &kBroker{
		addrs: cAddrs,
		codec: DefaultMarshaler{},
		opts:  options,
	}
}

type subscriber struct {
	kBroker       *kBroker
	consumerGroup sarama.ConsumerGroup
	topic         string
	opts          broker.SubscribeOptions
}

type publication struct {
	topic         string
	err           error
	consumerGroup sarama.ConsumerGroup
	kafkaMessage  *sarama.ConsumerMessage
	brokerMessage *broker.Message
	session       sarama.ConsumerGroupSession
	timestamp     time.Time
}

func (p *publication) Topic() string {
	return p.topic
}

func (p *publication) Message() *broker.Message {
	return p.brokerMessage
}

func (p *publication) Ack() error {
	p.session.MarkMessage(p.kafkaMessage, "")
	return nil
}

func (p *publication) Error() error {
	return p.err
}

func (p *publication) Timestamp() time.Time {
	return p.timestamp
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Unsubscribe() error {
	if err := s.consumerGroup.Close(); err != nil {
		return err
	}

	k := s.kBroker
	k.scMutex.Lock()
	defer k.scMutex.Unlock()

	for i, cg := range k.consumerGroups {
		if cg == s.consumerGroup {
			k.consumerGroups = append(k.consumerGroups[:i], k.consumerGroups[i+1:]...)
			return nil
		}
	}

	return nil
}

func (k *kBroker) Address() string {
	if len(k.addrs) > 0 {
		return k.addrs[0]
	}
	return DefaultKafkaBroker
}

func (k *kBroker) Connect() error {
	if k.connected {
		return nil
	}

	k.scMutex.Lock()
	if k.client != nil {
		k.scMutex.Unlock()
		return nil
	}
	k.scMutex.Unlock()

	pconfig := k.getProducerConfig(context.Background())
	// For implementation reasons, the SyncProducer requires
	// `Producer.Return.Errors` and `Producer.Return.Successes`
	// to be set to true in its configuration.
	pconfig.Producer.Return.Successes = true
	pconfig.Producer.Return.Errors = true

	c, err := sarama.NewClient(k.addrs, pconfig)
	if err != nil {
		return err
	}

	var (
		ap                   sarama.AsyncProducer
		p                    sarama.SyncProducer
		errChan, successChan = k.getAsyncProduceChan()
	)

	// Because error chan must require, so only error chan
	// If set the error chan, will use async produce
	// else use sync produce
	// only keep one client resource, is c variable
	if errChan != nil {
		ap, err = sarama.NewAsyncProducerFromClient(c)

		// opentelemetry tracing
		ap := otelsarama.WrapAsyncProducer(pconfig, ap)

		if err != nil {
			return err
		}
		// When the ap closed, the Errors() & Successes() channel will be closed
		// So the goroutine will auto exit
		go func() {
			for v := range ap.Errors() {
				errChan <- v
			}
		}()

		if successChan != nil {
			go func() {
				for v := range ap.Successes() {
					successChan <- v
				}
			}()
		}
	} else {
		p, err = sarama.NewSyncProducerFromClient(c)

		// opentelemetry tracing
		p = otelsarama.WrapSyncProducer(pconfig, p)

		if err != nil {
			return err
		}
	}

	k.scMutex.Lock()
	k.client = c
	if p != nil {
		k.syncProducer = p
	}
	if ap != nil {
		k.asyncProducer = ap
	}
	k.consumerGroups = make([]sarama.ConsumerGroup, 0)
	k.connected = true

	// request-reply pattern
	k.resps = sync.Map{}
	k.respSubscribers = sync.Map{}

	k.scMutex.Unlock()

	return nil
}

func (k *kBroker) Disconnect() error {
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	for _, consumer := range k.consumerGroups {
		consumer.Close()
	}
	k.consumerGroups = nil
	if k.syncProducer != nil {
		k.syncProducer.Close()
	}
	if k.asyncProducer != nil {
		k.asyncProducer.Close()
	}
	if err := k.client.Close(); err != nil {
		return err
	}
	k.connected = false

	// request-reply pattern
	k.resps = sync.Map{}
	k.respSubscribers = sync.Map{}

	return nil
}

func (k *kBroker) Init(opts ...broker.BrokerOption) error {
	for _, o := range opts {
		o(&k.opts)
	}
	var cAddrs []string
	for _, addr := range k.opts.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{DefaultKafkaBroker}
	}
	k.addrs = cAddrs
	return nil
}

func (k *kBroker) Options() broker.BrokerOptions {
	return k.opts
}

func (k *kBroker) Publish(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	// not used here
	_ = broker.NewPublishOptions(opts...)

	return k.sendMessage(ctx, topic, msg)
}

func (k *kBroker) PublishAndReceive(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) (*broker.Message, error) {
	options := broker.PublishOptions{
		ReplyToTopic: fmt.Sprintf("%s.reply", topic),
		Timeout:      RequestReplyTimeout,
	}

	for _, opt := range opts {
		opt(&options)
	}

	var (
		replyTopic         = options.ReplyToTopic
		replyConsumerGroup = options.ReplyConsumerGroup
		timeout            = options.Timeout
		errChan            = make(chan error)
		msgChan            = make(chan *broker.Message, 1)
	)

	err := k.sendMessage(ctx, topic, msg)
	if err != nil {
		return nil, err
	}

	// Create a channel to receive reply messages
	correlationId, correlationIdOk := msg.Headers[CorrelationIdHeader]
	if !correlationIdOk {
		return nil, fmt.Errorf("missing correlation id in message")
	}
	k.resps.Store(correlationId, msgChan)

	// Subscribe for reply topic if didn't
	go func() {
		if _, ok := k.respSubscribers.LoadOrStore(replyTopic, true); !ok {
			var subOpts = make([]broker.SubscribeOption, 0)
			if len(replyConsumerGroup) != 0 {
				subOpts = append(subOpts, broker.WithSubscribeGroup(replyConsumerGroup))
			}

			_, err := k.Subscribe(replyTopic, func(ctx context.Context, e broker.Event) error {
				go func() {
					if e.Message() == nil {
						return
					}

					cId, correlationIdOk := e.Message().Headers[CorrelationIdHeader]
					if !correlationIdOk {
						return
					}

					msgChan, msgChanOk := k.resps.LoadAndDelete(cId)
					if msgChanOk {
						msgChan.(chan *broker.Message) <- e.Message()
					}
				}()
				return nil
			}, subOpts...)

			if err != nil {
				errChan <- err
				k.respSubscribers.Delete(replyTopic)
			}
		}
	}()

	select {
	case body := <-msgChan:
		return body, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(timeout):
		// remove processed channel
		k.resps.Delete(correlationId)
		return nil, broker.RequestTimeoutResponse{
			Timeout: timeout,
		}
	}
}

func (k *kBroker) sendMessage(ctx context.Context, topic string, msg *broker.Message) error {
	kMsg, err := k.codec.Marshal(topic, msg)
	if err != nil {
		return fmt.Errorf("failed to marshal to kafka message: %w", err)
	}

	// opentelemetry tracing
	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(kMsg))

	if k.asyncProducer != nil {
		k.asyncProducer.Input() <- kMsg
		return nil
	} else if k.syncProducer != nil {
		_, _, err := k.syncProducer.SendMessage(kMsg)
		return err
	}
	return errors.New(`no connection resources available`)
}

func (k *kBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	start := time.Now()

	opt := broker.NewSubscribeOptions(opts...)

	// we need to create a new client per consumer
	cg, err := k.getSaramaConsumerGroup(opt.Context, opt.Group)
	if err != nil {
		return nil, err
	}

	csHandler := &consumerGroupHandler{
		handler: handler,
		subopts: opt,
		kopts:   k.opts,
		cg:      cg,
		ready:   make(chan bool),
		codec:   k.codec,
	}

	// Wrap instrumentation
	otelCsHandler := otelsarama.WrapConsumerGroupHandler(csHandler)

	ctx := context.Background()
	topics := []string{topic}
	go func() {
		for {
			select {
			case err := <-cg.Errors():
				if err != nil {
					k.log(ctx, logger.ErrorLevel, "consumer error: %s", err)
				}
			default:
				err := cg.Consume(ctx, topics, otelCsHandler)
				switch err {
				case sarama.ErrClosedConsumerGroup:
					return
				case nil:
					csHandler.ready = make(chan bool)
					continue
				default:
					k.log(ctx, logger.ErrorLevel, "consumer error: %s", err)
				}
			}
		}
	}()

	// wait until consumer group running
	<-csHandler.ready

	k.log(ctx, logger.InfoLevel, "Subcribed to topic: %s. Consumer group: %s. Duration: %dms", topic, opt.Group, time.Since(start).Milliseconds())

	return &subscriber{
		kBroker:       k,
		consumerGroup: cg,
		opts:          opt,
		topic:         topic,
	}, nil
}

func (k *kBroker) getProducerConfig(pContext context.Context) *sarama.Config {
	config := DefaultProducerConfig
	if c, ok := k.opts.Context.Value(producerConfigKey{}).(*sarama.Config); ok {
		config = c
	}
	return config
}

func (k *kBroker) getConsumerConfig(cContext context.Context) *sarama.Config {
	config := DefaultConsumerConfig
	if c, ok := k.opts.Context.Value(consumerConfigKey{}).(*sarama.Config); ok {
		config = c
	}

	// the oldest supported version is V0_10_2_0
	if !config.Version.IsAtLeast(sarama.V0_10_2_0) {
		config.Version = sarama.V0_10_2_0
	}

	if offset, ok := cContext.Value(initialOffsetKey{}).(int64); ok {
		config.Consumer.Offsets.Initial = offset
	}

	config.Consumer.Return.Errors = true

	return config
}

func (k *kBroker) getSaramaConsumerGroup(cContext context.Context, groupID string) (sarama.ConsumerGroup, error) {
	config := k.getConsumerConfig(cContext)
	cg, err := sarama.NewConsumerGroup(k.addrs, groupID, config)
	if err != nil {
		return nil, err
	}
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	k.consumerGroups = append(k.consumerGroups, cg)
	return cg, nil
}

func (k *kBroker) getAsyncProduceChan() (chan<- *sarama.ProducerError, chan<- *sarama.ProducerMessage) {
	var (
		errors    chan<- *sarama.ProducerError
		successes chan<- *sarama.ProducerMessage
	)
	if c, ok := k.opts.Context.Value(asyncProduceErrorKey{}).(chan<- *sarama.ProducerError); ok {
		errors = c
	}
	if c, ok := k.opts.Context.Value(asyncProduceSuccessKey{}).(chan<- *sarama.ProducerMessage); ok {
		successes = c
	}
	return errors, successes
}

func (k *kBroker) String() string {
	return "kafka broker implementation"
}

func (k *kBroker) log(ctx context.Context, level logger.Level, message string, args ...interface{}) {
	k.getLogger().Logf(ctx, level, message, args...)
}

func (k *kBroker) getLogger() logger.Logger {
	if k.opts.Logger != nil {
		return k.opts.Logger
	}
	return logger.DefaultLogger
}
