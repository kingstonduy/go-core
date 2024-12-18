package dispatcher

import (
	"context"

	"github.com/kingstonduy/go-core/transport"
	"github.com/kingstonduy/go-core/transport/broker"
)

type CommandHandlerFunc func(ctx context.Context, esCommand transport.Command) error

type DispatcherHandler interface {
	When(ctx context.Context, event transport.Command, broker broker.Broker) error
	RegisterHandler(commandType string, handler CommandHandlerFunc)
}
