#### Usage

```go

func GetKafkaBroker(cfg config.Configure, logger logger.Logger) broker.Broker {
	var (
		addrs             = strings.Split(cfg.GetString("KAFKA_BROKERS"), ",")
		TLSEnabled        = cfg.GetBool("KAFKA_TLS_ENABLED")
		TLSSkipVerify     = cfg.GetBool("KAFKA_TLS_SKIP_VERIFY")
		TLSCaCertFile     = cfg.GetString("KAFKA_CA_CERT_FILE")
		TLSClientKeyFile  = cfg.GetString("KAFKA_CLIENT_KEY_FILE")
		TLSClientCertFile = cfg.GetString("KAFKA_CLIENT_CERT_FILE")
		SASLEnabled       = cfg.GetBool("KAFKA_SASL_ENABLED")
		SASLAlgorithm     = cfg.GetString("KAFKA_SASL_ALGORITHM")
		SASLUser          = cfg.GetString("KAFKA_SASL_USER")
		SASLPassword      = cfg.GetString("KAFKA_SASL_PASSWORD")
	)
	var config = &kafka.KafkaBrokerConfig{
		Addresses:         addrs,
		TLSEnabled:        TLSEnabled,
		TLSSkipVerify:     TLSSkipVerify,
		TLSCaCertFile:     TLSCaCertFile,
		TLSClientCertFile: TLSClientCertFile,
		TLSClientKeyFile:  TLSClientKeyFile,
		SASLEnabled:       SASLEnabled,
		SASLAlgorithm:     SASLAlgorithm,
		SASLUser:          SASLUser,
		SASLPassword:      SASLPassword,
	}

	br, err := kafka.GetKafkaBroker(
		config,
		broker.WithLogger(logger),
	)

	ctx := context.Background()

	if err != nil {
		logger.Error(ctx, "Failted to create kafka broker")
		panic(err)
	}

	if err := br.Connect(); err != nil {
		logger.Error(ctx, "Failted to connect to kafka broker")
		panic(err)
	} else {
		logger.Info(ctx, "Connected to kafka broker")
	}

	return br
}

```


```go
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
		err = kBroker.Publish(context.Background(), replyTopic, &broker.Message{
			Headers: msg.Headers,
			Body:    resultByte,
		})

		if err != nil {
            // TODO
		}

		return nil
	}, broker.WithSubscribeGroup("benchmark.test"))
```

```go
    msg, err := kBroker.PublishAndReceive(context.Background(), requestTopic, &broker.Message{
            Headers: map[string]string{
                "correlationId": uuid.New().String(),
            },
            Body: reqByte,
        }, broker.WithPublishReplyToTopic(replyTopic),
            broker.WithReplyConsumerGroup("benchmark.test"))

```

```go
    err = kBroker.Publish(context.Background(), replyTopic, &broker.Message{
                Headers: msg.Headers,
                Body:    resultByte,
            })

```