package credis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/cache"
	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/logger"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type redisCacheClient struct {
	rClient redis.UniversalClient
	opts    cache.Options
}

func NewRedisClient(opts ...cache.Option) (cache.CacheClient, error) {
	options := cache.NewOptions(opts...)

	rClient, err := newUniversalClient(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}

	return &redisCacheClient{
		opts:    options,
		rClient: rClient,
	}, nil
}

func newUniversalClient(options cache.Options) (redis.UniversalClient, error) {
	if options.Context == nil {
		options.Context = context.Background()
	}

	opts, ok := options.Context.Value(redisOptionsContextKey{}).(redis.UniversalOptions)

	var rClient redis.UniversalClient
	if !ok {
		addr := cache.DefaultAddress
		if len(options.Address) > 0 {
			addr = options.Address
		}

		redisOptions, err := redis.ParseURL(addr)
		if err != nil {
			redisOptions = &redis.Options{Addr: addr}
		}

		rClient = redis.NewClient(redisOptions)
		if err := redisotel.InstrumentTracing(rClient); err != nil {
			return nil, err
		}

		return rClient, nil
	}

	if len(opts.Addrs) == 0 && len(options.Address) > 0 {
		opts.Addrs = []string{options.Address}
	}

	rClient = redis.NewUniversalClient(&opts)
	if err := redisotel.InstrumentTracing(rClient); err != nil {
		return nil, err
	}

	return rClient, nil
}

// SetNX implements cache.CacheClient.
func (r *redisCacheClient) SetNX(ctx context.Context, key string, values interface{}, expiration time.Duration) *redis.BoolCmd {
	if expiration == 0 {
		// default expiration
		expiration = r.opts.Expiration
	}

	bytes, err := json.Marshal(values)
	if err != nil {
		logger.Errorf(ctx, "failed to marshal value: %w", err)
	}
	return r.rClient.SetNX(ctx, key, bytes, expiration)
}

// Del implements cache.CacheClient.
func (r *redisCacheClient) Del(ctx context.Context, keys ...string) error {
	return r.rClient.Del(ctx, keys...).Err()
}

// Expire implements cache.CacheClient.
func (r *redisCacheClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.rClient.Expire(ctx, key, expiration).Err()
}

// FlushAll implements cache.CacheClient.
func (r *redisCacheClient) FlushAll(ctx context.Context) error {
	return r.rClient.FlushAll(ctx).Err()
}

// Get implements cache.CacheClient.
func (r *redisCacheClient) Get(ctx context.Context, key string, dest interface{}) (time.Duration, error) {
	val, err := r.rClient.Get(ctx, key).Bytes()
	if err != nil && err == redis.Nil {
		return 0, cache.ErrKeyNotFound
	} else if err != nil {
		return 0, err
	}

	dur, err := r.rClient.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(val, &dest)
	return dur, err
}

// SAdd implements cache.CacheClient.
func (r *redisCacheClient) SAdd(ctx context.Context, key string, members ...any) error {
	return r.rClient.SAdd(ctx, key, members...).Err()
}

// SMembers implements cache.CacheClient.
func (r *redisCacheClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.rClient.SMembers(ctx, key).Result()
}

// Set implements cache.CacheClient.
func (r *redisCacheClient) Set(ctx context.Context, key string, values any, expiration time.Duration) error {
	if expiration == 0 {
		// default expiration
		expiration = r.opts.Expiration
	}

	bytes, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.rClient.Set(ctx, key, bytes, expiration).Err()
}

// TTL implements cache.CacheClient.
func (r *redisCacheClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.rClient.TTL(ctx, key).Result()
}

// String implements cache.CacheClient.
func (r *redisCacheClient) String() string {
	return "redis"
}

func (r *redisCacheClient) log(ctx context.Context, level logger.Level, message string, args ...interface{}) {
	if !r.logEnabled() {
		return
	}
	logger := r.opts.Logger
	logger.Logf(ctx, level, message, args...)
}

func (r *redisCacheClient) logEnabled() bool {
	return r.opts.Logger != nil
}
