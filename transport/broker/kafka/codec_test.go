package kafka

import (
	"testing"

	"github.com/kingstonduy/go-core/transport/broker"

	"github.com/IBM/sarama"
)

func GetCodec() Codec {
	return DefaultMarshaler{}
}

func TestMarshaler(t *testing.T) {
	topic := "topic"
	var message = broker.Message{
		Headers: map[string]string{
			"CorrelationId": "123",
			"Key":           "value",
		},
		Body: []byte("value"),
	}

	var c = GetCodec()
	kMsg, err := c.Marshal(topic, &message)
	if err != nil {
		t.Error(err)
	}

	for _, h := range kMsg.Headers {
		v, ok := message.Headers[string(h.Key)]
		if !ok {
			t.Errorf("Failed to get message header: %s", string(h.Key))
		}

		if v != string(h.Value) {
			t.Errorf("Expected %v, got %v", v, string(h.Value))
		}
	}

	kBody, err := kMsg.Value.Encode()
	if err != nil {
		t.Error(err)
	}
	if string(message.Body) != string(kBody) {
		t.Errorf("Expected body: %v, got %v", string(message.Body), string(kBody))
	}

	if topic != kMsg.Topic {
		t.Errorf("Expected topic: %v, got %v", topic, kMsg.Topic)
	}
}

func TestUnmarshaler(t *testing.T) {
	var c = GetCodec()
	kMsg := sarama.ConsumerMessage{
		Topic: "topic",
		Headers: []*sarama.RecordHeader{
			{
				Key:   []byte("CorrelationId"),
				Value: []byte("123"),
			},
			{
				Key:   []byte("key"),
				Value: []byte("value"),
			},
		},
		Value: []byte("value"),
	}

	message, err := c.Unmarshal(&kMsg)
	if err != nil {
		t.Error(err)
	}

	for _, h := range kMsg.Headers {
		v, ok := message.Headers[string(h.Key)]
		if !ok {
			t.Errorf("Failed to get message header: %s", string(h.Key))
		}

		if v != string(h.Value) {
			t.Errorf("Expected %v, got %v", v, string(h.Value))
		}
	}

	if string(message.Body) != string(kMsg.Value) {
		t.Errorf("Expected body: %v, got %v", string(message.Body), string(kMsg.Value))
	}

}

func TestMarshalerWithKey(t *testing.T) {
	topic := "topic"
	var message = broker.Message{
		Key: []byte("message_key"),
		Headers: map[string]string{
			"CorrelationId": "123",
			"Key":           "value",
		},
		Body: []byte("value"),
	}

	var c = GetCodec()
	kMsg, err := c.Marshal(topic, &message)
	if err != nil {
		t.Error(err)
	}

	for _, h := range kMsg.Headers {
		v, ok := message.Headers[string(h.Key)]
		if !ok {
			t.Errorf("Failed to get message header: %s", string(h.Key))
		}

		if v != string(h.Value) {
			t.Errorf("Expected %v, got %v", v, string(h.Value))
		}
	}

	kBody, err := kMsg.Value.Encode()
	if err != nil {
		t.Error(err)
	}

	key, err := kMsg.Key.Encode()
	if err != nil {
		t.Error(err)
	}

	if string(message.Key) != string(key) {
		t.Errorf("Expected body: %v, got %v", string(message.Key), string(key))
	}

	if string(message.Body) != string(kBody) {
		t.Errorf("Expected body: %v, got %v", string(message.Body), string(kBody))
	}

	if topic != kMsg.Topic {
		t.Errorf("Expected topic: %v, got %v", topic, kMsg.Topic)
	}

}

func TestUnmarshalerWithKey(t *testing.T) {
	var c = GetCodec()
	kMsg := sarama.ConsumerMessage{
		Key:   []byte("message_key"),
		Topic: "topic",
		Headers: []*sarama.RecordHeader{
			{
				Key:   []byte("CorrelationId"),
				Value: []byte("123"),
			},
			{
				Key:   []byte("key"),
				Value: []byte("value"),
			},
		},
		Value: []byte("value"),
	}

	message, err := c.Unmarshal(&kMsg)
	if err != nil {
		t.Error(err)
	}

	for _, h := range kMsg.Headers {
		v, ok := message.Headers[string(h.Key)]
		if !ok {
			t.Errorf("Failed to get message header: %s", string(h.Key))
		}

		if v != string(h.Value) {
			t.Errorf("Expected %v, got %v", v, string(h.Value))
		}
	}

	if string(message.Body) != string(kMsg.Value) {
		t.Errorf("Expected body: %v, got %v", string(message.Body), string(kMsg.Value))
	}

	if string(message.Key) != string(kMsg.Key) {
		t.Errorf("Expected Key: %v, got %v", string(message.Key), string(kMsg.Key))
	}

}
