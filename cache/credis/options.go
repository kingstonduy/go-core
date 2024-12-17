package credis

import (
	"context"

	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/cache"
	"github.com/redis/go-redis/v9"
)

type redisOptionsContextKey struct{}

// WithRedisOptions sets advanced options for redis.
func WithRedisOptions(options redis.UniversalOptions) cache.Option {
	return func(o *cache.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, redisOptionsContextKey{}, options)
	}
}
