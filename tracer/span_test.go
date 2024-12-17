package trace

import (
	"context"
	"fmt"
	"testing"
)

func TestGetSpanInfo(t *testing.T) {
	ctx := context.Background()

	spanInfo := GetSpanInfo(ctx)

	fmt.Printf("spanInfo: %+v\n", spanInfo)

	// enrich spanInfo
	spanInfo = SpanInfo{
		TraceID:       "traceID",
		SpanID:        "spanID",
		ServiceDomain: "serviceDomain",
		OperatorName:  "operatorName",
		StepName:      "stepName",
		ClientID:      "clientID",
		SystemID:      "systemID",
		From:          "from",
		To:            "to",
	}

	ctx = InjectSpanInfo(ctx, spanInfo)

	spanInfo = GetSpanInfo(ctx)

	fmt.Printf("spanInfo: %+v\n", spanInfo)
}
