package metrics

import "testing"

func TestDefaultMetrics(t *testing.T) {
	IncrCounterWithLabels([]string{}, 3.4, []Label{})
}
