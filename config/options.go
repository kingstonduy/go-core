package config

import (
	"context"
)

type ConfigOptions struct {
	// configuration file
	ConfigFile string

	// default configuration
	Defaults map[string]interface{}

	// mapping the environment variables from outside
	AutomaticEnv bool

	// mapping the environment variables to struct fields based on tag
	TagName string

	// configuration context
	Context context.Context
}

type ConfigOption func(*ConfigOptions)

// configuration file. Default: .env
func WithConfigFile(configFile string) ConfigOption {
	return func(o *ConfigOptions) {
		o.ConfigFile = configFile
	}
}

// default configurations
func WithDefaults(defaults map[string]interface{}) ConfigOption {
	return func(o *ConfigOptions) {
		if o.Defaults == nil {
			o.Defaults = make(map[string]interface{})
		}

		for k, v := range defaults {
			o.Defaults[k] = v
		}
	}
}

// default configuration
func WithDefault(key string, value interface{}) ConfigOption {
	return func(o *ConfigOptions) {
		if o.Defaults == nil {
			o.Defaults = make(map[string]interface{})
		}

		o.Defaults[key] = value
	}
}

func WithContext(ctx context.Context) ConfigOption {
	return func(o *ConfigOptions) {
		o.Context = ctx
	}
}

// mapping the environment variables to struct fields based on tag. Default: mapstructure
func WithTagName(tag string) ConfigOption {
	return func(o *ConfigOptions) {
		o.TagName = tag
	}
}

// mapping the environment variables from outside. Default: true
func WithAutomaticEnv(auto bool) ConfigOption {
	return func(o *ConfigOptions) {
		o.AutomaticEnv = auto
	}
}

// NewOptions returns a new options struct.
func NewOptions(opts ...ConfigOption) ConfigOptions {
	options := ConfigOptions{
		ConfigFile:   DefaultConfigFile, // no expiration
		Defaults:     make(map[string]interface{}),
		AutomaticEnv: true,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
