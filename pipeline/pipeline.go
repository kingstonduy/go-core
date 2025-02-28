package pipeline

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/ahmetb/go-linq/v3"
	"github.com/kingstonduy/go-core/logger"
)

func init() {
	// Register default behaviors
	err := RegisterRequestPipelineBehaviors(
		NewTracingBehavior(),
		NewRequestLoggingBehavior(),
		NewMetricsBehavior(),
		NewErrorHandlingBehavior(),
		NewValidationBehavior(),
	)

	if err != nil {
		logger.Error(context.TODO(), "failed to register default behaviors: %s", err)
	}
}

// RequestHandlerFunc is a continuation for the next task to execute in the pipeline
type RequestHandlerFunc func(ctx context.Context) (interface{}, error)

// PipelineBehavior is a Pipeline behavior for wrapping the inner handler.
type PipelineBehavior interface {
	Handle(ctx context.Context, request interface{}, next RequestHandlerFunc) (interface{}, error)
}

type RequestHandler[TRequest any, TResponse any] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}

type RequestHandlerFactory[TRequest any, TResponse any] func() RequestHandler[TRequest, TResponse]

type NotificationHandler[TNotification any] interface {
	Handle(ctx context.Context, notification TNotification) error
}

type NotificationHandlerFactory[TNotification any] func() NotificationHandler[TNotification]

var requestHandlersRegistrations = map[reflect.Type]interface{}{}
var notificationHandlersRegistrations = map[reflect.Type][]interface{}{}
var pipelineBehaviours []interface{} = []interface{}{}

type Unit struct{}

func registerRequestHandler[TRequest any, TResponse any](handler any) error {
	var request TRequest
	requestType := reflect.TypeOf(request)

	_, exist := requestHandlersRegistrations[requestType]
	if exist {
		// each request in request/response strategy should have just one handler
		return fmt.Errorf("registered handler already exists in the registry for message %s", requestType.String())
	}

	requestHandlersRegistrations[requestType] = handler

	return nil
}

// RegisterRequestHandler register the request handler to mediatr registry.
func RegisterRequestHandler[TRequest any, TResponse any](handler RequestHandler[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest, TResponse](handler)
}

// RegisterRequestHandlerFactory register the request handler factory to mediatr registry.
func RegisterRequestHandlerFactory[TRequest any, TResponse any](factory RequestHandlerFactory[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest, TResponse](factory)
}

// RegisterRequestPipelineBehaviors register the request behaviors to mediatr registry.
func RegisterRequestPipelineBehaviors(behaviours ...PipelineBehavior) error {
	for _, behavior := range behaviours {
		behaviorType := reflect.TypeOf(behavior)

		existsPipe := existsPipeType(behaviorType)
		if existsPipe {
			return fmt.Errorf("registered behavior already exists in the registry.")
		}

		pipelineBehaviours = append(pipelineBehaviours, behavior)
	}

	return nil
}

func registerNotificationHandler[TEvent any](handler any) error {
	var event TEvent
	eventType := reflect.TypeOf(event)

	handlers, exist := notificationHandlersRegistrations[eventType]
	if !exist {
		notificationHandlersRegistrations[eventType] = []interface{}{handler}
		return nil
	}

	notificationHandlersRegistrations[eventType] = append(handlers, handler)

	return nil
}

// RegisterNotificationHandler register the notification handler to mediatr registry.
func RegisterNotificationHandler[TEvent any](handler NotificationHandler[TEvent]) error {
	return registerNotificationHandler[TEvent](handler)
}

// RegisterNotificationHandlerFactory register the notification handler factory to mediatr registry.
func RegisterNotificationHandlerFactory[TEvent any](factory NotificationHandlerFactory[TEvent]) error {
	return registerNotificationHandler[TEvent](factory)
}

// RegisterNotificationHandlers register the notification handlers to mediatr registry.
func RegisterNotificationHandlers[TEvent any](handlers ...NotificationHandler[TEvent]) error {
	if len(handlers) == 0 {
		return errors.New("no handlers provided")
	}

	for _, handler := range handlers {
		err := RegisterNotificationHandler(handler)
		if err != nil {
			return err
		}
	}

	return nil
}

// RegisterNotificationHandlersFactories register the notification handlers factories to mediatr registry.
func RegisterNotificationHandlersFactories[TEvent any](factories ...NotificationHandlerFactory[TEvent]) error {
	if len(factories) == 0 {
		return errors.New("no handlers provided")
	}

	for _, handler := range factories {
		err := RegisterNotificationHandlerFactory[TEvent](handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func ClearRequestRegistrations() {
	requestHandlersRegistrations = map[reflect.Type]interface{}{}
}

func ClearNotificationRegistrations() {
	notificationHandlersRegistrations = map[reflect.Type][]interface{}{}
}

func buildRequestHandler[TRequest any, TResponse any](handler any) (RequestHandler[TRequest, TResponse], bool) {
	handlerValue, ok := handler.(RequestHandler[TRequest, TResponse])
	if !ok {
		factory, ok := handler.(RequestHandlerFactory[TRequest, TResponse])
		if !ok {
			return nil, false
		}

		return factory(), true
	}

	return handlerValue, true
}

// Send the request to its corresponding request handler.
func Send[TRequest any, TResponse any](ctx context.Context, request TRequest) (TResponse, error) {
	requestType := reflect.TypeOf(request)

	var response TResponse
	handler, ok := requestHandlersRegistrations[requestType]
	if !ok {
		// request-response strategy should have exactly one handler and if we can't find a corresponding handler, we should return an error
		return *new(TResponse), fmt.Errorf("no handler for request %T", request)
	}

	handlerValue, ok := buildRequestHandler[TRequest, TResponse](handler)
	if !ok {
		return *new(TResponse), fmt.Errorf("handler for request %T is not a Handler", request)
	}

	if len(pipelineBehaviours) > 0 {
		var reversPipes = reversOrder(pipelineBehaviours)

		var lastHandler RequestHandlerFunc = func(ctx context.Context) (interface{}, error) {
			return handlerValue.Handle(ctx, request)
		}

		aggregateResult := linq.From(reversPipes).AggregateWithSeedT(lastHandler, func(next RequestHandlerFunc, pipe PipelineBehavior) RequestHandlerFunc {
			pipeValue := pipe
			nexValue := next

			var handlerFunc RequestHandlerFunc = func(ctx context.Context) (interface{}, error) {
				return pipeValue.Handle(ctx, request, nexValue)
			}

			return handlerFunc
		})

		v := aggregateResult.(RequestHandlerFunc)
		resp, err := v(ctx)

		if resp != nil {
			response = resp.(TResponse)
		}

		if err != nil {
			return response, err
		}

	} else {
		res, err := handlerValue.Handle(ctx, request)
		if err != nil {
			return res, err
		}

		response = res
	}

	return response, nil
}

func buildNotificationHandler[TNotification any](handler any) (NotificationHandler[TNotification], bool) {
	handlerValue, ok := handler.(NotificationHandler[TNotification])
	if !ok {
		factory, ok := handler.(NotificationHandlerFactory[TNotification])
		if !ok {
			return nil, false
		}

		return factory(), true
	}

	return handlerValue, true
}

// Publish the notification event to its corresponding notification handler.
func Publish[TNotification any](ctx context.Context, notification TNotification) error {
	eventType := reflect.TypeOf(notification)

	handlers, ok := notificationHandlersRegistrations[eventType]
	if !ok {
		// notification strategy should have zero or more handlers, so it should run without any error if we can't find a corresponding handler
		return nil
	}

	for _, handler := range handlers {
		handlerValue, ok := buildNotificationHandler[TNotification](handler)

		if !ok {
			return fmt.Errorf("handler for notification %T is not a Handler", notification)
		}

		err := handlerValue.Handle(ctx, notification)
		if err != nil {
			return err
		}
	}

	return nil
}

func reversOrder(values []interface{}) []interface{} {
	var reverseValues []interface{}

	for i := len(values) - 1; i >= 0; i-- {
		reverseValues = append(reverseValues, values[i])
	}

	return reverseValues
}

func existsPipeType(p reflect.Type) bool {
	for _, pipe := range pipelineBehaviours {
		if reflect.TypeOf(pipe) == p {
			return true
		}
	}

	return false
}
