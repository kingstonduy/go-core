package credis

import (
	"context"
	"testing"
	"time"

	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/cache"
	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/trace"
	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/trace/otel"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func getCacheClient() cache.CacheClient {
	client, err := NewRedisClient(
		WithRedisOptions(
			redis.UniversalOptions{
				// MasterName: "redis-master",
				Addrs: []string{"127.0.0.1:6379"},
			},
		),
	)
	if err != nil {
		panic(err)
	}

	return client
}

type DataTest struct {
	Data string
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func TestSetKey(t *testing.T) {
	tests := []KeyValue{
		{
			Key:   "test1",
			Value: "value",
		},
		{
			Key:   "test2",
			Value: 1,
		},
		{
			Key:   "test3",
			Value: true,
		},
		{
			Key: "test4",
			Value: DataTest{
				Data: "data-test",
			},
		},
		{
			Key:   "test5",
			Value: []byte("string"),
		},
	}

	client := getCacheClient()
	for _, d := range tests {
		err := client.Set(context.Background(), d.Key, d.Value, time.Duration(0))
		if err != nil {
			t.Error(err)
		}
	}

	// get data
	var result0 string
	key0 := tests[0].Key
	expected0 := tests[0].Value
	_, err := client.Get(context.Background(), key0, &result0)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected0, result0)

	var result1 int
	key1 := tests[1].Key
	expected1 := tests[1].Value
	_, err = client.Get(context.Background(), key1, &result1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected1, result1)

	var result2 bool
	key2 := tests[2].Key
	expected2 := tests[2].Value
	_, err = client.Get(context.Background(), key2, &result2)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected2, result2)

	var result3 DataTest
	key3 := tests[3].Key
	expected3 := tests[3].Value
	_, err = client.Get(context.Background(), key3, &result3)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected3, result3)

	var result4 []byte
	key4 := tests[4].Key
	expected4 := tests[4].Value
	_, err = client.Get(context.Background(), key4, &result4)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected4, result4)

}

func TestDeleteKey(t *testing.T) {
	key := "test"
	value := DataTest{
		Data: "data-test",
	}
	client := getCacheClient()
	err := client.Set(context.Background(), key, value, time.Second*20)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.Del(context.Background(), key)
	if err != nil {
		t.Error(err)
	}

	var data interface{}
	_, err = client.Get(context.Background(), key, data)

	assert.Contains(t, err.Error(), "key not found in cache")

}

func TestSetExpireKey(t *testing.T) {
	key := "test"
	value := DataTest{
		Data: "data-test",
	}
	expire := time.Second * 10

	client := getCacheClient()
	err := client.Set(context.Background(), key, value, time.Duration(0))
	if err != nil {
		t.Error(err)
		return
	}

	err = client.Expire(context.Background(), key, expire)
	if err != nil {
		t.Error(err)
	}

	dur, err := client.TTL(context.Background(), "test")
	if err != nil {
		t.Error(err)
	}

	assert.EqualValues(t, dur, expire)
}

func TestFlushAll(t *testing.T) {
	key := "test"
	value := DataTest{
		Data: "data-test",
	}

	client := getCacheClient()
	err := client.Set(context.Background(), key, value, time.Duration(0))
	if err != nil {
		t.Error(err)
		return
	}

	err = client.FlushAll(context.Background())
	if err != nil {
		t.Error(err)
	}

	var data interface{}
	_, err = client.Get(context.Background(), key, data)
	assert.EqualValues(t, err, cache.ErrKeyNotFound)
}

func TestSAdd(t *testing.T) {
	key := "test"
	value1 := "test1"

	value2 := "test2"

	ctx := context.Background()

	client := getCacheClient()
	err := client.Del(ctx, key)
	if err != nil {
		t.Error(err)
	}
	err = client.SAdd(ctx, key, value1, value2)
	if err != nil {
		t.Error(err)
		return
	}

	data, err := client.SMembers(ctx, key)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, len(data), 2)
}

func TestWithTracing(t *testing.T) {
	tracer, err := otel.NewOpenTelemetryTracer(
		context.Background(),
		trace.WithTraceServiceName("cache-service"),
		trace.WithServiceVersion("1.0.4"),
		trace.WithTraceExporterEndpoint("localhost:4318"),
	)

	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()
	ctx, f := tracer.StartTracing(ctx, "start-cache-tracing")

	tests := []KeyValue{
		{
			Key:   "test1",
			Value: "value",
		},
		{
			Key:   "test2",
			Value: 1,
		},
		{
			Key:   "test3",
			Value: true,
		},
		{
			Key: "test4",
			Value: DataTest{
				Data: "data-test",
			},
		},
		{
			Key:   "test5",
			Value: []byte("string"),
		},
	}

	client := getCacheClient()
	for _, d := range tests {
		err := client.Set(ctx, d.Key, d.Value, time.Duration(0))
		if err != nil {
			t.Error(err)
		}
	}

	// get data
	var result0 string
	key0 := tests[0].Key
	expected0 := tests[0].Value
	_, err = client.Get(ctx, key0, &result0)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected0, result0)

	var result1 int
	key1 := tests[1].Key
	expected1 := tests[1].Value
	_, err = client.Get(ctx, key1, &result1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected1, result1)

	var result2 bool
	key2 := tests[2].Key
	expected2 := tests[2].Value
	_, err = client.Get(ctx, key2, &result2)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected2, result2)

	var result3 DataTest
	key3 := tests[3].Key
	expected3 := tests[3].Value
	_, err = client.Get(ctx, key3, &result3)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected3, result3)

	var result4 []byte
	key4 := tests[4].Key
	expected4 := tests[4].Value
	_, err = client.Get(ctx, key4, &result4)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected4, result4)

	f(ctx)
	time.Sleep(5 * time.Second)
}

func TestDefaultCacheClient(t *testing.T) {
	tests := []KeyValue{
		{
			Key:   "test1",
			Value: "value",
		},
		{
			Key:   "test2",
			Value: 1,
		},
		{
			Key:   "test3",
			Value: true,
		},
		{
			Key: "test4",
			Value: DataTest{
				Data: "data-test",
			},
		},
		{
			Key:   "test5",
			Value: []byte("string"),
		},
	}

	client := getCacheClient()
	cache.SetDefaultCacheClient(client)
	for _, d := range tests {
		err := cache.Set(context.Background(), d.Key, d.Value, 10*time.Second)
		if err != nil {
			t.Error(err)
		}
	}

	assert.Equal(t, client, cache.DefaultCacheClient)

}
