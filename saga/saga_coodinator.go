package saga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gammazero/workerpool"
	"github.com/kingstonduy/go-core/logger"
	"github.com/kingstonduy/go-core/transport/broker"
)

var ErrAbortSaga = errors.New("saga aborted")

type SagaStartOptions struct {
	brokerSubcription []broker.SubscribeOption
	sagaRoutes        []Saga
}

type SagaStartOption func(*SagaStartOptions)

// The request header mapped to response headers
func WithStartSubcription(subcriptionOptions []broker.SubscribeOption) SagaStartOption {
	return func(options *SagaStartOptions) {
		options.brokerSubcription = subcriptionOptions
	}
}

// The request header mapped to response headers
func WithSagaRoute(sagaRoutes []Saga) SagaStartOption {
	return func(options *SagaStartOptions) {
		options.sagaRoutes = sagaRoutes
	}
}

type SEC struct {
	Broker     broker.Broker
	Sagas      map[string]Saga
	WorkerPool *workerpool.WorkerPool
	SagaTopic  string
	quit       chan struct{}
}

func NewSec(broker broker.Broker,
	workerpool *workerpool.WorkerPool,
	sagaTopic string) *SEC {
	return &SEC{
		Broker:     broker,
		WorkerPool: workerpool,
		SagaTopic:  sagaTopic,
		quit:       make(chan struct{}),
		Sagas:      make(map[string]Saga),
	}
}

// RegisterSaga adds a Saga to the SEC.
func (s *SEC) RegisterSaga(saga Saga) {
	s.Sagas[saga.name] = saga
}

func (s *SEC) Start(ctx context.Context, opts ...SagaStartOption) error {
	options := SagaStartOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	if len(options.sagaRoutes) == 0 {
		logger.Errorf(ctx, "no saga route for handler")
		return fmt.Errorf("no saga route for handler")
	}

	for _, saga := range options.sagaRoutes {
		s.RegisterSaga(saga)
	}

	subscriber, err := s.Broker.Subscribe(s.SagaTopic, func(c context.Context, e broker.Event) error {
		// prevent closure
		if err := e.Ack(); err != nil {
			return err
		}

		ctx := c
		event := e
		s.WorkerPool.Submit(func() {
			if event.Message() == nil || len(event.Message().Body) == 0 {
				logger.Errorf(ctx, "Topic: %s. Empty messsage body", event.Topic())
			}
			var sagaCommand SagaCommand
			// keyAggregateId := string(event.Message().Key)
			if err := json.Unmarshal(event.Message().Body, &sagaCommand); err != nil {
				logger.Errorf(ctx, "Topic: %s. Can Not Unmarshal messsage body", event.Topic())
				return
			}

			err := s.ProcessCommand(ctx, sagaCommand)
			if err != nil {
				logger.Errorf(ctx, "Failed to handle message - Error: %s", err.Error())
			}
		})
		return nil
	}, options.brokerSubcription...)
	if err != nil {
		logger.Errorf(ctx, "Failed to consume to topic: %s", s.SagaTopic)
		subscriber.Unsubscribe() //nolint
		return err
	}
	go func() {
		defer subscriber.Unsubscribe() //nolint
		<-s.quit
	}()
	return nil
}

func (s *SEC) ProcessCommand(ctx context.Context, sagaCommand SagaCommand) error {
	saga, ok := s.Sagas[sagaCommand.SagaName]
	if !ok {
		return fmt.Errorf("no saga with name %s exists", sagaCommand.SagaName)
	}

	switch sagaCommand.Name {
	case BeginSagaCommand:
		nextTransaction := saga.FirstTransaction()
		if nextTransaction == "" {
			return s.Write(ctx, EndSaga(sagaCommand.SagaName, sagaCommand.SagaID))
		}

		return s.Write(ctx, BeginTransaction(sagaCommand.SagaName, sagaCommand.SagaID, nextTransaction, sagaCommand.SagaParams))
	case BeginTransactionCommand:
		execErr := saga.ExecuteTransaction(ctx, sagaCommand.TransactionID, sagaCommand.SagaParams)
		if execErr != nil {
			if errors.Is(execErr, ErrAbortSaga) {
				// abort saga, need to compensate transactions to the save point
				return s.Write(ctx, AbortSaga(sagaCommand.SagaName, sagaCommand.SagaID, sagaCommand.TransactionID))
			}
			// abort transaction, need to repeat this transaction again
			return s.Write(ctx, AbortTransaction(sagaCommand.SagaName, sagaCommand.SagaID, sagaCommand.TransactionID, sagaCommand.SagaParams))
		}

		return s.Write(ctx, EndTransaction(sagaCommand.SagaName, sagaCommand.SagaID, sagaCommand.TransactionID, sagaCommand.SagaParams))
	case AbortTransactionCommand:
		return s.Write(ctx, BeginTransaction(sagaCommand.SagaName, sagaCommand.SagaID, sagaCommand.TransactionID, sagaCommand.SagaParams))
	case AbortSagaCommand:
		return s.Write(ctx, EndTransactionCompensate(sagaCommand.SagaName, sagaCommand.SagaID, sagaCommand.TransactionID, saga.Compensation(sagaCommand.TransactionID), sagaCommand.SagaParams))
	case EndTransactionCommand:
		nextTransaction := saga.Next(sagaCommand.TransactionID)
		if sagaCommand.CompensationID != "" {
			nextTransaction = sagaCommand.CompensationID
		}

		if nextTransaction == "" {
			return s.Write(ctx, EndSaga(sagaCommand.SagaName, sagaCommand.SagaID))
		}

		return s.Write(ctx, BeginTransaction(sagaCommand.SagaName, sagaCommand.SagaID, nextTransaction, sagaCommand.SagaParams))
	case EndSagaCommand:
		return nil
	default:
		return fmt.Errorf("unknow command %d", sagaCommand.Name)
	}
}

// RegisterSaga adds a Saga to the SEC.
func (s *SEC) Write(ctx context.Context, sagaCommand SagaCommand) error {
	v, err := json.Marshal(sagaCommand)
	if err != nil {
		logger.Errorf(ctx, "failed command marshalling: %v", err)
		return fmt.Errorf("failed command marshalling: %v", err)
	}
	msg := broker.Message{
		Key:  []byte(sagaCommand.SagaID),
		Body: v,
	}
	err = s.Broker.Publish(ctx, s.SagaTopic, &msg)
	if err != nil {
		logger.Errorf(ctx, "Produce saga command failed: %v", err)
		return fmt.Errorf("produce saga command failed: %v", err)
	}
	return nil
}

// RegisterSaga adds a Saga to the SEC.
func (s *SEC) Stop(ctx context.Context) {
	close(s.quit)
}
