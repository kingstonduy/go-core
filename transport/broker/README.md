```go
func WithT24AdapterSubscription() BrokerServerStartOption {
	return func(b *BrokerServer) error {
		_, err := b.broker.Subscribe(b.topics.RequestT24Adapter, func(c context.Context, e broker.Event) error {
			// prevent closure
			ctx := c
			event := e
			b.workerpool.Submit(func() {
				if event.Message() == nil && len(event.Message().Body) == 0 {
					logger.Infof(ctx, "Topic: %s. Empty message body", event.Topic())
				}

				broker.HandleBrokerEvent[*t24.T24AdapterRequest, *t24.T24AdapterResponse](
					ctx,
					event,
					broker.WithOnRequestHandledFunc(func(ctx context.Context, res interface{}) {
						// dosomething
					}),
				)
			})
			return nil
		}, b.GetSubscriptionOptions()...)
		if err != nil {
			return err
		}

		return nil
	}
}
```