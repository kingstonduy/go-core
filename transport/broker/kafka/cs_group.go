package kafka

import (
	"context"

	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/transport/broker"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"go.opentelemetry.io/otel"
)

// consumerGroupHandler is the implementation of sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler broker.Handler
	subopts broker.SubscribeOptions
	kopts   broker.BrokerOptions
	cg      sarama.ConsumerGroup
	ready   chan bool
	codec   Codec
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			// opentelemetry tracing
			ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(msg))

			if !ok {
				h.log(ctx, logger.InfoLevel, "[kafka consumer] message channel was closed")
				return nil
			}

			if msg == nil {
				continue
			}

			m, err := h.codec.Unmarshal(msg)
			if err != nil {
				h.log(ctx, logger.ErrorLevel, "[kafka consumer]: failed to unmarshal consumed message: %v", err)
				continue
			}

			p := &publication{brokerMessage: m, topic: msg.Topic, kafkaMessage: msg, consumerGroup: h.cg, session: session, timestamp: msg.Timestamp}
			logger.Fields(
				map[string]interface{}{
					logger.FIELD_OPERATOR_NAME: p.topic,
					logger.FIELD_STEP_NAME:     "message-received",
				},
			).Info(ctx, broker.MakeStringLogsKafka(ctx, *p.brokerMessage))
			err = h.handler(ctx, p)
			if err == nil && h.subopts.AutoAck {
				session.MarkMessage(msg, "")
			} else if err != nil {
				p.err = err
				errHandler := h.kopts.ErrorHandler
				if errHandler != nil {
					errHandler(ctx, p) //nolint
				} else {
					h.log(ctx, logger.ErrorLevel, "[kafka] subscriber error: %v", err)
				}
			}
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *consumerGroupHandler) log(ctx context.Context, level logger.Level, message string, args ...interface{}) {
	c.getLogger().Logf(ctx, level, message, args...)
}

func (c *consumerGroupHandler) getLogger() logger.Logger {
	if c.kopts.Logger != nil {
		return c.kopts.Logger
	}

	return logger.DefaultLogger
}
