package rabbitmq_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/transport/broker"
	"github.com/kingstonduy/go-core/transport/broker/rabbitmq"
	"github.com/stretchr/testify/assert"
)

// import (
// 	"context"
// 	"encoding/json"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/kingstonduy/go-core/logger"
// 	"github.com/kingstonduy/go-core/transport/broker"
// 	"github.com/kingstonduy/go-core/transport/broker/rabbitmq"
// 	"go-micro.dev/v4"
// 	"go-micro.dev/v4/server"
// )

type Example struct{}

func init() {
	rabbitmq.DefaultRabbitURL = "amqp://rabbit:rabbit@127.0.0.1:5672"
}

type TestEvent struct {
	Name string    `json:"name"`
	Age  int       `json:"age"`
	Time time.Time `json:"time"`
}

func (e *Example) Handler(ctx context.Context, r interface{}) error {
	return nil
}

type RabbitMQTest struct {
	T *testing.T
}

func TestRabbitMq(t *testing.T) {
	t.Run("1. Test broker", func(t *testing.T) {
		test := RabbitMQTest{t}
		// test.TestCreateBroker()
		// test.TestConnectBroker()
		// test.TestConnectWithoutExchange()
		test.TestConnectWithDirectExchange()
	})
}

func (r RabbitMQTest) TestCreateBroker() {
	t := r.T
	b := rabbitmq.NewBroker()
	assert.NotNil(t, b)
	assert.Equal(t, b.Options().Logger, logger.DefaultLogger)
}

func (r RabbitMQTest) TestConnectBroker() {
	t := r.T
	b := rabbitmq.NewBroker()
	if err := b.Connect(); err != nil {
		t.Logf("cant conect to broker, skip: %v", err)
		t.Error(err)
	}
}

func (r RabbitMQTest) TestConnectWithoutExchange() {
	t := r.T
	b := rabbitmq.NewBroker(rabbitmq.WithoutExchange())
	wg := new(sync.WaitGroup)

	wg.Add(1)
	if err := b.Connect(); err != nil {
		t.Logf("cant conect to broker, skip: %v", err)
		t.Error(err)
	}

	subOpts := []broker.SubscribeOption{
		broker.WithSubscribeQueue("direct.queue"),
		broker.WithSubscribeAutoAck(false),
	}

	_, err := b.Subscribe("direct.queue", func(ctx context.Context, e broker.Event) error {
		body := e.Message().Body
		var res TestEvent
		err := json.Unmarshal(body, &res)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Receive event: %v", string(body))
		assert.Equal(t, res.Age, 16)
		assert.Equal(t, res.Name, "test")

		wg.Done()
		return nil
	}, subOpts...)

	if err != nil {
		t.Error(err)
	}

	go func(t *testing.T) {
		time.Sleep(5 * time.Second)
		logger.Infof(context.Background(), "pub event")

		jsonData, _ := json.Marshal(&TestEvent{
			Name: "test",
			Age:  16,
		})

		err := b.Publish(context.Background(), "direct.queue", &broker.Message{
			Body: jsonData,
		},
			rabbitmq.DeliveryMode(2),
			rabbitmq.ContentType("application/json"),
		)

		if err != nil {
			t.Error(err)
		}
	}(t)

	wg.Wait()
}

func (r RabbitMQTest) TestConnectWithDirectExchange() {
	t := r.T
	b := rabbitmq.NewBroker(
		rabbitmq.ExchangeType(rabbitmq.ExchangeTypeDirect),
		rabbitmq.ExchangeName("direct.test"),
	)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	if err := b.Connect(); err != nil {
		t.Logf("cant connect to broker, skip: %v", err)
		t.Error(err)
	}

	subOpts := []broker.SubscribeOption{
		broker.WithSubscribeQueue("direct.exchange.queue"),
		broker.WithSubscribeAutoAck(false),
	}

	_, err := b.Subscribe("direct.exchange.queue", func(ctx context.Context, e broker.Event) error {
		body := e.Message().Body
		var res TestEvent
		err := json.Unmarshal(body, &res)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Receive event: %v", string(body))
		assert.Equal(t, res.Age, 16)
		assert.Equal(t, res.Name, "test")

		wg.Done()
		return nil
	}, subOpts...)

	if err != nil {
		t.Error(err)
	}

	go func(t *testing.T) {
		time.Sleep(5 * time.Second)
		logger.Infof(context.Background(), "pub event")

		jsonData, _ := json.Marshal(&TestEvent{
			Name: "test",
			Age:  16,
		})

		err := b.Publish(context.Background(), "direct.exchange.queue", &broker.Message{
			Body: jsonData,
		},
			rabbitmq.DeliveryMode(2),
			rabbitmq.ContentType("application/json"),
		)

		if err != nil {
			t.Error(err)
		}
	}(t)

	wg.Wait()
}

// func TestWithoutExchange(t *testing.T) {

// 	b := rabbitmq.NewBroker(rabbitmq.WithoutExchange())
// 	b.Init()
// 	if err := b.Connect(); err != nil {
// 		t.Logf("cant conect to broker, skip: %v", err)
// 		t.Skip()
// 	}

// 	s := server.NewServer(server.Broker(b))

// 	service := micro.NewService(
// 		micro.Server(s),
// 		micro.Broker(b),
// 	)
// 	brkrSub := broker.NewSubscribeOptions(
// 		broker.Queue("direct.queue"),
// 		broker.DisableAutoAck(),
// 		rabbitmq.DurableQueue(),
// 	)
// 	// Register a subscriber
// 	err := micro.RegisterSubscriber(
// 		"direct.queue",
// 		service.Server(),
// 		func(ctx context.Context, evt *TestEvent) error {
// 			logger.Logf(logger.InfoLevel, "receive event: %+v", evt)
// 			return nil
// 		},
// 		server.SubscriberContext(brkrSub.Context),
// 		server.SubscriberQueue("direct.queue"),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	go func() {
// 		time.Sleep(5 * time.Second)
// 		logger.Logf(logger.InfoLevel, "pub event")
// 		jsonData, _ := json.Marshal(&TestEvent{
// 			Name: "test",
// 			Age:  16,
// 		})
// 		err := b.Publish("direct.queue", &broker.Message{
// 			Body: jsonData,
// 		},
// 			rabbitmq.DeliveryMode(2),
// 			rabbitmq.ContentType("application/json"))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}()

// 	// service.Init()

// 	if err := service.Run(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestFanoutExchange(t *testing.T) {
// 	b := rabbitmq.NewBroker(rabbitmq.ExchangeType(rabbitmq.ExchangeTypeFanout), rabbitmq.ExchangeName("fanout.test"))
// 	b.Init()
// 	if err := b.Connect(); err != nil {
// 		t.Logf("cant conect to broker, skip: %v", err)
// 		t.Skip()
// 	}

// 	s := server.NewServer(server.Broker(b))

// 	service := micro.NewService(
// 		micro.Server(s),
// 		micro.Broker(b),
// 	)
// 	brkrSub := broker.NewSubscribeOptions(
// 		broker.Queue("fanout.queue"),
// 		broker.DisableAutoAck(),
// 		rabbitmq.DurableQueue(),
// 	)
// 	// Register a subscriber
// 	err := micro.RegisterSubscriber(
// 		"fanout.queue",
// 		service.Server(),
// 		func(ctx context.Context, evt *TestEvent) error {
// 			logger.Logf(logger.InfoLevel, "receive event: %+v", evt)
// 			return nil
// 		},
// 		server.SubscriberContext(brkrSub.Context),
// 		server.SubscriberQueue("fanout.queue"),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	go func() {
// 		time.Sleep(5 * time.Second)
// 		logger.Logf(logger.InfoLevel, "pub event")
// 		jsonData, _ := json.Marshal(&TestEvent{
// 			Name: "test",
// 			Age:  16,
// 		})
// 		err := b.Publish("fanout.queue", &broker.Message{
// 			Body: jsonData,
// 		},
// 			rabbitmq.DeliveryMode(2),
// 			rabbitmq.ContentType("application/json"))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}()

// 	// service.Init()

// 	if err := service.Run(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestDirectExchange(t *testing.T) {
// 	b := rabbitmq.NewBroker(rabbitmq.ExchangeType(rabbitmq.ExchangeTypeDirect), rabbitmq.ExchangeName("direct.test"))
// 	b.Init()
// 	if err := b.Connect(); err != nil {
// 		t.Logf("cant conect to broker, skip: %v", err)
// 		t.Skip()
// 	}

// 	s := server.NewServer(server.Broker(b))

// 	service := micro.NewService(
// 		micro.Server(s),
// 		micro.Broker(b),
// 	)
// 	brkrSub := broker.NewSubscribeOptions(
// 		broker.Queue("direct.exchange.queue"),
// 		broker.DisableAutoAck(),
// 		rabbitmq.DurableQueue(),
// 	)
// 	// Register a subscriber
// 	err := micro.RegisterSubscriber(
// 		"direct.exchange.queue",
// 		service.Server(),
// 		func(ctx context.Context, evt *TestEvent) error {
// 			logger.Logf(logger.InfoLevel, "receive event: %+v", evt)
// 			return nil
// 		},
// 		server.SubscriberContext(brkrSub.Context),
// 		server.SubscriberQueue("direct.exchange.queue"),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	go func() {
// 		time.Sleep(5 * time.Second)
// 		logger.Logf(logger.InfoLevel, "pub event")
// 		jsonData, _ := json.Marshal(&TestEvent{
// 			Name: "test",
// 			Age:  16,
// 		})
// 		err := b.Publish("direct.exchange.queue", &broker.Message{
// 			Body: jsonData,
// 		},
// 			rabbitmq.DeliveryMode(2),
// 			rabbitmq.ContentType("application/json"))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}()

// 	// service.Init()

// 	if err := service.Run(); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestTopicExchange(t *testing.T) {
// 	b := rabbitmq.NewBroker()
// 	b.Init()
// 	if err := b.Connect(); err != nil {
// 		t.Logf("cant conect to broker, skip: %v", err)
// 		t.Skip()
// 	}

// 	s := server.NewServer(server.Broker(b))

// 	service := micro.NewService(
// 		micro.Server(s),
// 		micro.Broker(b),
// 	)
// 	brkrSub := broker.NewSubscribeOptions(
// 		broker.Queue("topic.exchange.queue"),
// 		broker.DisableAutoAck(),
// 		rabbitmq.DurableQueue(),
// 	)
// 	// Register a subscriber
// 	err := micro.RegisterSubscriber(
// 		"my-test-topic",
// 		service.Server(),
// 		func(ctx context.Context, evt *TestEvent) error {
// 			logger.Logf(logger.InfoLevel, "receive event: %+v", evt)
// 			return nil
// 		},
// 		server.SubscriberContext(brkrSub.Context),
// 		server.SubscriberQueue("topic.exchange.queue"),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	go func() {
// 		time.Sleep(5 * time.Second)
// 		logger.Logf(logger.InfoLevel, "pub event")
// 		jsonData, _ := json.Marshal(&TestEvent{
// 			Name: "test",
// 			Age:  16,
// 		})
// 		err := b.Publish("my-test-topic", &broker.Message{
// 			Body: jsonData,
// 		},
// 			rabbitmq.DeliveryMode(2),
// 			rabbitmq.ContentType("application/json"))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}()

// 	// service.Init()

// 	if err := service.Run(); err != nil {
// 		t.Fatal(err)
// 	}
// }
