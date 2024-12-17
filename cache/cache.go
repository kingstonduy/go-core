package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	DefaultCacheClient = newNoopsCacheClient()

	DefaultExpiration time.Duration = 0

	DefaultAddress = "127.0.0.1:6379"

	ErrItemExpired error = errors.New("item has expired")

	ErrKeyNotFound error = errors.New("key not found in cache")
)

func SetDefaultCacheClient(c CacheClient) {
	DefaultCacheClient = c
}

type CacheClient interface {
	Get(ctx context.Context, key string, dest interface{}) (time.Duration, error)
	// Duration: -2 if the key does not exist
	// Duration: -1 if the key exists but has no associated expire
	TTL(ctx context.Context, key string) (time.Duration, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Set(ctx context.Context, key string, values interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	FlushAll(ctx context.Context) error
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	String() string
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
}

func Get(ctx context.Context, key string, dest interface{}) (time.Duration, error) {
	return DefaultCacheClient.Get(ctx, key, dest)
}

// Duration: -2 if the key does not exist
// Duration: -1 if the key exists but has no associated expire
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return DefaultCacheClient.TTL(ctx, key)
}

func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return DefaultCacheClient.Expire(ctx, key, expiration)
}

func Set(ctx context.Context, key string, values interface{}, expiration time.Duration) error {
	return DefaultCacheClient.Set(ctx, key, values, expiration)
}

func Del(ctx context.Context, keys ...string) error {
	return DefaultCacheClient.Del(ctx, keys...)
}

func FlushAll(ctx context.Context) error {
	return DefaultCacheClient.FlushAll(ctx)
}

func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return DefaultCacheClient.SAdd(ctx, key, members...)
}

func SMembers(ctx context.Context, key string) ([]string, error) {
	return DefaultCacheClient.SMembers(ctx, key)
}

func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return DefaultCacheClient.SetNX(ctx, key, value, expiration)
}
