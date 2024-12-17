package cache

import (
	"context"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type noopsCacheClient struct {
}

// Del implements CacheClient.
func (n *noopsCacheClient) Del(ctx context.Context, keys ...string) error {
	n.noopsWarning()
	return nil
}

// Expire implements CacheClient.
func (n *noopsCacheClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	n.noopsWarning()
	return nil
}

// FlushAll implements CacheClient.
func (n *noopsCacheClient) FlushAll(ctx context.Context) error {
	n.noopsWarning()
	return nil
}

// Get implements CacheClient.
func (n *noopsCacheClient) Get(ctx context.Context, key string, dest interface{}) (time.Duration, error) {
	n.noopsWarning()
	return 0, nil
}

// SAdd implements CacheClient.
func (n *noopsCacheClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	n.noopsWarning()
	return nil
}

// SMembers implements CacheClient.
func (n *noopsCacheClient) SMembers(ctx context.Context, key string) ([]string, error) {
	n.noopsWarning()
	return nil, nil
}

// Set implements CacheClient.
func (n *noopsCacheClient) Set(ctx context.Context, key string, values interface{}, expiration time.Duration) error {
	n.noopsWarning()
	return nil
}

// String implements CacheClient.
func (n *noopsCacheClient) String() string {
	n.noopsWarning()
	return ""
}

// TTL implements CacheClient.
func (n *noopsCacheClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	n.noopsWarning()
	return 0, nil
}

func (n *noopsCacheClient) noopsWarning() {
	log.Print("[WARN] No default cache client was set. Using noops cache client as default. Set the default cache client to do all functions\n")
}

// SetNX implements CacheClient.
func (n *noopsCacheClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	n.noopsWarning()
	return nil
}

func newNoopsCacheClient() CacheClient {
	return &noopsCacheClient{}
}
