package kafka

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kingstonduy/go-core/transport/broker"
	"github.com/stretchr/testify/assert"
)

func TestInitialOffset(t *testing.T) {
	var config = &KafkaBrokerConfig{
		Addresses: []string{"localhost:9092"},
	}

	br, err := GetKafkaBroker(
		config,
	)

	if err != nil {
		t.Fatal(err)
	}

	options := []broker.SubscribeOption{
		InitialOffset(-2),
		broker.WithSubscribeGroup("test.offset"),
		broker.WithSubscribeAutoAck(false),
	}

	_, err = br.Subscribe(
		"test.offset",
		func(ctx context.Context, e broker.Event) error {
			if e.Message() != nil {
				t.Logf("Received message: %v", e.Message().Body)
			}

			if err := e.Ack(); err != nil {
				return err
			}

			return nil
		}, options...)

	if err != nil {
		t.Error(err)
	}

	time.Sleep(10 * time.Second)
}

func TestMessageKey(t *testing.T) {
	topic := "test_message_key"
	wg := sync.WaitGroup{}
	config := &KafkaBrokerConfig{
		Addresses: []string{"127.0.0.1:9092"},
	}

	br, err := GetKafkaBroker(
		config,
	)

	if err := br.Connect(); err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Fatal(err)
	}

	msg := broker.Message{
		Key:  []byte(uuid.New().String()),
		Body: []byte("message_body"),
		Headers: map[string]string{
			"header_key": "header_value",
		},
	}

	wg.Add(1)
	_, err = br.Subscribe(topic, func(ctx context.Context, e broker.Event) error {
		if message := e.Message(); message != nil {
			assert.Equal(t, string(message.Key), string(msg.Key))
			assert.Equal(t, string(message.Body), string(message.Body))
			for k, v := range message.Headers {
				assert.Equal(t, msg.Headers[k], v)
			}
			wg.Done()
		}
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	err = br.Publish(context.TODO(), topic, &msg)
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}
