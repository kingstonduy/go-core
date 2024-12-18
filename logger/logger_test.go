package logger

import (
	"context"
	"sync"
	"testing"
)

func TestNoopsLogger(t *testing.T) {
	ctx := context.Background()
	Logf(ctx, InfoLevel, "logging: %s", "data")
	Tracef(ctx, "logging: %s", "data")
	Infof(ctx, "logging: %s", "data")
	Warnf(ctx, "logging: %s", "data")
	Debugf(ctx, "logging: %s", "data")
	Errorf(ctx, "logging: %s", "data")
	Fatalf(ctx, "logging: %s", "data")
}

func TestConcurrency(t *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(100)
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		go func() {
			Infof(ctx, "logging: %s", "data")
			wg.Done()
		}()
	}
	wg.Wait()
}
