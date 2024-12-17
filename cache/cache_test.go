package cache

import (
	"context"
	"testing"
	"time"
)

type DataTest struct {
	Data string
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func TestDefault(t *testing.T) {
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

	for _, d := range tests {
		err := Set(context.Background(), d.Key, d.Value, time.Duration(0))
		if err != nil {
			t.Error(err)
		}
	}
}
