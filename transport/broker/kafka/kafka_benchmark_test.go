package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/trace"
	"github.com/kingstonduy/go-core/trace/otel"
	"github.com/kingstonduy/go-core/transport/broker"

	"github.com/google/uuid"
)

func getKafkaBroker() broker.Broker {
	var config = &KafkaBrokerConfig{
		Addresses: []string{"localhost:9092"},
	}

	br, err := GetKafkaBroker(
		config,
	)

	if err != nil {
		panic(err)
	}

	return br
}

func getTracer() trace.Tracer {
	var (
		serviceName = "test-service"
		exporter    = "localhost:4318"
	)

	tracer, err := otel.NewOpenTelemetryTracer(
		context.Background(),
		trace.WithTraceServiceName(serviceName),
		trace.WithTraceExporterEndpoint(exporter),
	)

	if err != nil {
		panic(err)
	}

	return tracer
}

type KRequestType struct {
	Number int
}

type KResponseType struct {
	Result float64
}

func BenchmarkKafka(b *testing.B) {
	var (
		kBroker      = getKafkaBroker()
		tracer       = getTracer()
		requestTopic = "go.clean.test.benchmark.request"
		replyTopic   = "go.clean.test.benchmark.reply"
		errCount     int64
	)

	err := kBroker.Connect()
	if err != nil {
		b.Fail()
	}

	_, err = kBroker.Subscribe(requestTopic, func(ctx context.Context, e broker.Event) error {
		msg := e.Message()
		if msg == nil {
			return broker.EmptyMessageError{}
		}

		var req KRequestType
		err := json.Unmarshal(msg.Body, &req)
		if err != nil {
			return broker.InvalidDataFormatError{}
		}

		result := KResponseType{
			Result: math.Pow(float64(req.Number), 2),
		}

		resultByte, err := json.Marshal(result)
		if err != nil {
			return broker.InvalidDataFormatError{}
		}

		// pubish to response topic
		err = kBroker.Publish(ctx, replyTopic, &broker.Message{
			Headers: msg.Headers,
			Body:    resultByte,
		})

		if err != nil {
			b.Error(err)
			b.Fail()
		}

		return nil
	}, broker.WithSubscribeGroup("benchmark.test"))

	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	round := 5
	total := 10

	var totalTime int64
	var wg sync.WaitGroup
	wg.Add(total * round)
	for i := 0; i < round; i++ {
		go func(i int) {
			for j := 0; j < total; j++ {
				req := KRequestType{
					Number: rand.Intn(100),
				}
				ctx, f := tracer.StartTracing(context.Background(), fmt.Sprintf("operation-%v", i*j+1), trace.WithTraceRequest(req))
				reqByte, err := json.Marshal(req)
				if err != nil {
					b.Error(err)
				}

				start := time.Now()
				b.Logf("Start time: %v", start)
				msg, err := kBroker.PublishAndReceive(ctx, requestTopic, &broker.Message{
					Headers: map[string]string{
						"correlationId": uuid.New().String(),
					},
					Body: reqByte,
				}, broker.WithPublishReplyToTopic(replyTopic),
					broker.WithReplyConsumerGroup("benchmark.test"))

				if err != nil {
					errCount++
					b.Errorf("benchmark error: %v", err)
					f(ctx, trace.WithTraceErrorResponse(err))
				} else {
					var result map[string]float64
					json.Unmarshal(msg.Body, &result) //nolint
					expected := math.Pow(float64(req.Number), 2)
					if result["Result"] != expected {
						b.Errorf("Expected result: %v, got: %v", expected, result["Result"])
					}
					b.Logf("Message: %v. Expected result: %v, got: %v, Duration: %vms", string(msg.Headers[CorrelationIdHeader]), expected, result["Result"], time.Since(start).Milliseconds())
					f(ctx, trace.WithTraceResponse(result))
				}

				totalTime += time.Since(start).Milliseconds()
				wg.Done()
			}
		}(i)
	}

	wg.Wait()
	b.Logf("AVG: %vms", totalTime/int64(total*round))
	if err := kBroker.Disconnect(); err != nil {
		b.Error(err)
	}

	time.Sleep(20 * time.Second)

}

func TestPublishAndReceived(t *testing.T) {
	var (
		kBroker      = getKafkaBroker()
		tracer       = getTracer()
		requestTopic = "go.clean.test.benchmark.request"
		replyTopic   = "go.clean.test.benchmark.reply"
	)

	err := kBroker.Connect()
	if err != nil {
		t.Error(err)
	}

	_, err = kBroker.Subscribe(requestTopic, func(ctx context.Context, e broker.Event) error {
		msg := e.Message()
		if msg == nil {
			return broker.EmptyMessageError{}
		}

		var req KRequestType
		err := json.Unmarshal(msg.Body, &req)
		if err != nil {
			return broker.InvalidDataFormatError{}
		}

		result := KResponseType{
			Result: math.Pow(float64(req.Number), 2),
		}

		resultByte, err := json.Marshal(result)
		if err != nil {
			return broker.InvalidDataFormatError{}
		}

		// pubish to response topic
		err = kBroker.Publish(ctx, replyTopic, &broker.Message{
			Headers: msg.Headers,
			Body:    resultByte,
		})

		if err != nil {
			t.Error(err)
		}

		return nil
	}, broker.WithSubscribeGroup("benchmark.test"))

	if err != nil {
		t.Error(err)
	}

	req := KRequestType{
		Number: rand.Intn(100),
	}

	reqByte, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
	}

	ctx, f := tracer.StartTracing(context.Background(), "Publish and Receive", trace.WithTraceRequest(req))
	msg, err := kBroker.PublishAndReceive(ctx, requestTopic, &broker.Message{
		Body: reqByte,
	},
		broker.WithPublishReplyToTopic(replyTopic),
		broker.WithReplyConsumerGroup("benchmark.test"))

	if err != nil {
		t.Errorf("error: %v", err)
		f(ctx, trace.WithTraceErrorResponse(err))

	} else {
		t.Logf("%v", msg)
		f(ctx, trace.WithTraceResponse(msg))
	}

	time.Sleep(10 * time.Second)

}
