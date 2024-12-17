package viper

import (
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/kingstonduy/go-core/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	DefaultTagName = "mapstructure"
)

type ViperConfig struct {
	cfg        *viper.Viper
	configOpts config.ConfigOptions
}

func NewViperConfig(opts ...config.ConfigOption) (config.Configure, error) {
	options := config.NewOptions(opts...)

	vi := viper.New()

	// Read the environment variables
	if options.AutomaticEnv {
		vi.AutomaticEnv()
	}

	// Override environment file
	if value, ok := os.LookupEnv(config.EnvFileKey); ok && len(value) > 0 {
		options.ConfigFile = value
	}

	// Read the environment variables from the config file
	if len(options.ConfigFile) != 0 {
		if _, err := os.Stat(options.ConfigFile); err == nil {
			vi.SetConfigFile(options.ConfigFile)
			err := vi.ReadInConfig()
			if err != nil {
				return nil, err
			}
		}
	}

	// Set default configuration
	if len(options.Defaults) != 0 {
		for k, v := range options.Defaults {
			vi.SetDefault(k, v)
		}
	}

	return &ViperConfig{
		cfg:        vi,
		configOpts: options,
	}, nil
}

func (c *ViperConfig) bindEnvs(iface interface{}, parts ...string) {
	vi := c.cfg
	tagName := DefaultTagName
	if len(c.configOpts.TagName) != 0 {
		tagName = c.configOpts.TagName
	}

	ifv := reflect.ValueOf(iface)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}
	for i := 0; i < ifv.NumField(); i++ {
		v := ifv.Field(i)
		t := ifv.Type().Field(i)
		tv, ok := t.Tag.Lookup(tagName)
		if !ok {
			continue
		}
		if tv == ",squash" {
			c.bindEnvs(v.Interface(), parts...)
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			c.bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			vi.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

// Get implements ConfigInterface.
func (c *ViperConfig) Get(key string) interface{} {
	return c.cfg.Get(key)
}

// GetBool implements ConfigInterface.
func (c *ViperConfig) GetBool(key string) bool {
	return c.cfg.GetBool(key)
}

// GetDuration implements ConfigInterface.
func (c *ViperConfig) GetDuration(key string) time.Duration {
	return c.cfg.GetDuration(key)
}

// GetFloat64 implements ConfigInterface.
func (c *ViperConfig) GetFloat64(key string) float64 {
	return c.cfg.GetFloat64(key)
}

// GetInt implements ConfigInterface.
func (c *ViperConfig) GetInt(key string) int {
	return c.cfg.GetInt(key)
}

// GetInt32 implements ConfigInterface.
func (c *ViperConfig) GetInt32(key string) int32 {
	return c.cfg.GetInt32(key)
}

// GetInt64 implements ConfigInterface.
func (c *ViperConfig) GetInt64(key string) int64 {
	return c.cfg.GetInt64(key)
}

// GetIntSlice implements ConfigInterface.
func (c *ViperConfig) GetIntSlice(key string) []int {
	return c.cfg.GetIntSlice(key)
}

// GetString implements ConfigInterface.
func (c *ViperConfig) GetString(key string) string {
	return c.cfg.GetString(key)
}

// GetStringMap implements ConfigInterface.
func (c *ViperConfig) GetStringMap(key string) map[string]any {
	return c.cfg.GetStringMap(key)
}

// GetStringSlice implements ConfigInterface.
func (c *ViperConfig) GetStringSlice(key string) []string {
	return c.cfg.GetStringSlice(key)
}

// GetTime implements ConfigInterface.
func (c *ViperConfig) GetTime(key string) time.Time {
	return c.cfg.GetTime(key)
}

// GetUint implements ConfigInterface.
func (c *ViperConfig) GetUint(key string) uint {
	return c.cfg.GetUint(key)
}

// GetUint16 implements ConfigInterface.
func (c *ViperConfig) GetUint16(key string) uint16 {
	return c.cfg.GetUint16(key)
}

// GetUint32 implements ConfigInterface.
func (c *ViperConfig) GetUint32(key string) uint32 {
	return c.cfg.GetUint32(key)
}

// GetUint64 implements ConfigInterface.
func (c *ViperConfig) GetUint64(key string) uint64 {
	return c.cfg.GetUint64(key)
}

// Unmarshal implements config.Configure.
func (c *ViperConfig) Unmarshal(dest interface{}) error {
	c.bindEnvs(dest) // bind environment from outside
	return c.cfg.Unmarshal(
		dest,
		tagNameOption(c.configOpts.TagName),
		weakTypeOption(true),
	)
}

// UnmarshalKey implements config.Configure.
func (c *ViperConfig) UnmarshalKey(key string, dest interface{}) error {
	c.bindEnvs(dest) // bind environment from outside
	return c.cfg.UnmarshalKey(
		key,
		dest,
		tagNameOption(c.configOpts.TagName),
		weakTypeOption(true),
	)
}

func tagNameOption(tag string) viper.DecoderConfigOption {
	return func(cfg *mapstructure.DecoderConfig) {
		cfg.TagName = tag
	}
}
func weakTypeOption(b bool) viper.DecoderConfigOption {
	return func(cfg *mapstructure.DecoderConfig) {
		cfg.WeaklyTypedInput = b
	}
}
