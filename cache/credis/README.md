#### USAGE

```go

func GetRedisClient(cfg Configuration) cache.CacheClient {
	redisConfig := cfg.RedisConfig

	client, err := credis.NewRedisClient(
		credis.WithRedisOptions(
			redis.UniversalOptions{
				Addrs:    strings.Split(redisConfig.Addresses, ","),
				Username: redisConfig.Username,
				Password: redisConfig.Password,
			},
		),
		cache.WithDefaultExpiration(redisConfig.Expiration),
	)
	if err != nil {
		panic(err)
	}

	logger.Infof(context.TODO(), "Conntected to redis server: %s", redisConfig.Addresses)

	return client
}

```