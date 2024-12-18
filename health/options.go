package healthchecks

import "github.com/kingstonduy/go-core/logger"

type HealthCheckOptions struct {
	Name        string
	Version     string
	Description string
	Logger      logger.Logger
}

type HealthCheckOption func(*HealthCheckOptions)

func WithName(name string) HealthCheckOption {
	return func(options *HealthCheckOptions) {
		options.Name = name
	}
}

func WithDescription(description string) HealthCheckOption {
	return func(options *HealthCheckOptions) {
		options.Description = description
	}
}

func WithVersion(version string) HealthCheckOption {
	return func(options *HealthCheckOptions) {
		options.Version = version
	}
}

func WithLogger(logger logger.Logger) HealthCheckOption {
	return func(options *HealthCheckOptions) {
		options.Logger = logger
	}
}
