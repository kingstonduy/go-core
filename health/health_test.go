package healthchecks

import (
	"encoding/json"
	"testing"
	"time"
)

var heathCheckData = struct {
	name        string
	description string
	version     string
}{
	name:        "test-name",
	description: "test-description",
	version:     "test-version",
}

func TestGR(t *testing.T) {
	healthChecker := NewHealthChecker(
		WithName(heathCheckData.name),
		WithDescription(heathCheckData.description),
		WithVersion(heathCheckData.version),
	)

	healthChecker.AddLivenessCheck("test-gc", NewGoroutineChecker(10))

	detail := healthChecker.LivenessCheck()
	printDataDetail(t, detail)
}

func TestPing(t *testing.T) {
	healthChecker := NewHealthChecker(
		WithName(heathCheckData.name),
		WithDescription(heathCheckData.description),
		WithVersion(heathCheckData.version),
	)

	healthChecker.AddLivenessCheck("test-ping", NewPingChecker("http://google.com", "get", 200*time.Millisecond, nil, nil))

	detail := healthChecker.LivenessCheck()
	printDataDetail(t, detail)
}

func TestGC(t *testing.T) {
	healthChecker := NewHealthChecker(
		WithName(heathCheckData.name),
		WithDescription(heathCheckData.description),
		WithVersion(heathCheckData.version),
	)

	healthChecker.AddLivenessCheck("test-gc", NewGCMaxChecker(100*time.Millisecond))

	detail := healthChecker.LivenessCheck()
	printDataDetail(t, detail)
}

func printDataDetail(t *testing.T, detail ApplicationHealthDetailed) {
	dataJson, err := json.MarshalIndent(detail, "", "\t")
	if err != nil {
		t.Error(err)
	}

	t.Logf("Detail: %v", string(dataJson))
}
