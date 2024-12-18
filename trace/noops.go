package trace

import (
	"context"
	"log"
)

const (
	Warning = "No default tracer was set. Using noops tracer as default"
)

type noopsTrace struct{}

func newNoopsTrace() Tracer {
	return &noopsTrace{}
}

func (n *noopsTrace) ExtractSpanInfo(ctx context.Context) SpanInfo {
	n.noopsWarning()
	return SpanInfo{}
}

func (n *noopsTrace) StartTracing(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, SpanFinishFunc) {
	n.noopsWarning()
	return ctx, func(ctx context.Context, sfo ...SpanFinishOption) {}
}

func (n *noopsTrace) noopsWarning() {
	log.Print("[WARN] No default tracer was set. Using noops tracer as default. Set the default tracer to do all functions\n")
}
