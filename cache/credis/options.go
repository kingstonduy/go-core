package credis

import (
	"context"

	"github.com/kingstonduy/go-core/cache"
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
