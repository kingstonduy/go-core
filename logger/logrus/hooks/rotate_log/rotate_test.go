package rotateLog

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/logger/logrus"
	"github.com/kingstonduy/go-core/metadata"
)

func TestLogRotateStoreLogIntoFile(t *testing.T) {

	if err := os.Setenv(metadata.EnvApplicationName, "service-name"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv(metadata.EnvPodName, "pod-name"); err != nil {
		t.Fatal(err)
	}

	hook, err := NewRotateLogHook()

	if err != nil {
		t.Error(err)
	}

	log := logrus.NewLogrusLogger(
		logrus.WithHooks(hook),
	)

	log.Infof(context.TODO(), "Congratulations!")
	time.Sleep(1 * time.Second)

	log.Infof(context.TODO(), "Congratulations!")
	time.Sleep(2 * time.Second)
}

func TestLogEntryWritten(t *testing.T) {

	if err := os.Setenv(metadata.EnvApplicationName, "service-name-2"); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv(metadata.EnvPodName, "pod-name-2"); err != nil {
		t.Fatal(err)
	}

	var serviceName, podName string
	if value, ok := os.LookupEnv(metadata.EnvApplicationName); ok {
		serviceName = value
	}

	if value, ok := os.LookupEnv(metadata.EnvPodName); ok {
		podName = value
	}

	pattern := fmt.Sprintf("%s/%s/%s/%s.%s.log", metadata.BASE_PATH_LOG, serviceName, podName, podName, "%Y%m%d%H%M")

	hook, err := NewRotateLogHook(
		WithRotateLogFilePattern(pattern),
		WithRotateLogRotationTime(time.Minute),
	)
	if err != nil {
		t.Error(err)
	}

	log := logrus.NewLogrusLogger(
		logrus.WithHooks(hook),
	)

	log.Infof(context.TODO(), "Congratulations!")
	time.Sleep(1 * time.Second)

	log.Infof(context.TODO(), "Congratulations!")
	time.Sleep(2 * time.Second)
}
