package trace

import (
	"context"
	"testing"
)

func TestDefaultTracer(t *testing.T) {
	StartTracing(context.Background(), "default tracing")
	ExtractSpanInfo(context.Background())
}
