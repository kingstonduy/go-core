package viper

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kingstonduy/go-core/config"
	"github.com/stretchr/testify/assert"
)

var (
	SampleEnv = "sample.env"
)

func TestGetConfigure(t *testing.T) {
	c, err := NewViperConfig()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	assert.NotNil(t, c)
}

func TestGetConfigureWithEnvFileKey(t *testing.T) {
	os.Setenv(config.EnvFileKey, "sample2.env")
	c, err := NewViperConfig(
		config.WithConfigFile("sample.env"),
	)

	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.NotNil(t, c)
	assert.Equal(t, "sample2", c.GetString("CONFIG_SAMPLE"))
}

func TestGetConfigureWithDefaultValues(t *testing.T) {
	c, err := NewViperConfig(
		config.WithDefaults(map[string]interface{}{
			"DEFAULT_KEY": "DEFAULT_VALUE",
		}),
		config.WithDefault("DEFAULT_KEY_2", "DEFAULT_VALUE_2"),
	)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	assert.NotNil(t, c)
	assert.Equal(t, "DEFAULT_VALUE", c.GetString("DEFAULT_KEY"))
	assert.Equal(t, "DEFAULT_VALUE_2", c.GetString("DEFAULT_KEY_2"))

}

func TestGetConfigureWithAutomaticEnvTrue(t *testing.T) {
	c, err := NewViperConfig(
		config.WithDefaults(map[string]interface{}{
			"DEFAULT_KEY": "DEFAULT_VALUE",
		}),
		config.WithAutomaticEnv(true),
	)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	assert.NotNil(t, c)
	assert.NotNil(t, c.GetString("USERNAME"))
}

func TestGetConfigureWithAutomaticEnvFalse(t *testing.T) {
	c, err := NewViperConfig(
		config.WithDefaults(map[string]interface{}{
			"DEFAULT_KEY": "DEFAULT_VALUE",
		}),
		config.WithAutomaticEnv(false),
	)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	assert.NotNil(t, c)
	assert.Equal(t, "", c.GetString("USERNAME"))
}

func TestGetConfigureWithEnvFile(t *testing.T) {
	c, err := NewViperConfig(
		config.WithDefaults(map[string]interface{}{
			"DEFAULT_KEY": "DEFAULT_VALUE",
		}),
		config.WithConfigFile("sample.env"),
	)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	ti, err := time.Parse(time.RFC3339, "2020-08-19T15:24:39-07:00")
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.NotNil(t, c)
	assert.Equal(t, "config", c.Get("CONFIG"))
	assert.Equal(t, true, c.GetBool("CONFIG_BOOL"))
	assert.Equal(t, "string", c.GetString("CONFIG_STRING"))
	assert.Equal(t, 5, c.GetInt("CONFIG_NUMBER_INT"))
	assert.Equal(t, int32(-999999), c.GetInt32("CONFIG_NUMBER_INT32"))
	assert.Equal(t, int64(999999), c.GetInt64("CONFIG_NUMBER_INT64"))
	assert.Equal(t, uint(999999), c.GetUint("CONFIG_NUMBER_UINT"))
	assert.Equal(t, uint32(999999), c.GetUint32("CONFIG_NUMBER_UINT32"))
	assert.Equal(t, uint64(999999), c.GetUint64("CONFIG_NUMBER_UINT64"))
	assert.Equal(t, 999999.9, c.GetFloat64("CONFIG_NUMBER_FLOAT64"))
	assert.Equal(t, ti, c.GetTime("CONFIG_TIME"))
	assert.Equal(t, time.Second*5, c.GetDuration("CONFIG_TIME_DURATION"))
	assert.Equal(t, time.Second*5, c.GetDuration("CONFIG_TIME_DURATION"))
	// assert.Equal(t, []int{1, 2, 3}, c.GetIntSlice("CONFIG_INT_SLICE"))
	assert.Equal(t, []string{"str1", "str2", "str3"}, c.GetStringSlice("CONFIG_STRING_SLICE"))
	assert.Equal(t, map[string]interface{}{
		"str1": "str1",
		"str2": "str2",
		"str3": "str3",
	}, c.GetStringMap("CONFIG_STRING_MAP"))

}

type Config struct {
	IntField     int          `config:"CONFIG_INT_FIELD"`
	FloatField   float64      `config:"CONFIG_FLOAT_FIELD"`
	StringField  string       `config:"CONFIG_STRING_FIELD"`
	BoolField    bool         `config:"CONFIG_BOOL_FIELD"`
	NestedConfig NestedConfig `config:",squash"`
}

type NestedConfig struct {
	IntNestedField    int    `config:"CONFIG_INT_NESTED_FIELD"`
	StringNestedField string `config:"CONFIG_STRING_NESTED_FIELD"`
}

func TestGetConfigureUnmarshal(t *testing.T) {
	c, err := NewViperConfig(
		config.WithDefaults(map[string]interface{}{
			"DEFAULT_KEY": "DEFAULT_VALUE",
		}),
		config.WithTagName("config"),
		config.WithAutomaticEnv(true),
		config.WithConfigFile("sample.env"),
	)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	assert.NotNil(t, c)

	var config Config

	err = c.Unmarshal(&config)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 9, config.IntField)
	assert.Equal(t, float64(9.9), config.FloatField)
	assert.Equal(t, "sample", config.StringField)
	assert.Equal(t, true, config.BoolField)
	assert.Equal(t, 9, config.NestedConfig.IntNestedField)
	assert.Equal(t, "sample", config.NestedConfig.StringNestedField)

}

func TestGetBindEnv(t *testing.T) {
	c, err := NewViperConfig(
		config.WithTagName("config"),
	)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	assert.NotNil(t, c)

	// set OS Environment variables
	os.Setenv("CONFIG_INT_FIELD", "1")
	os.Setenv("CONFIG_FLOAT_FIELD", "2")
	os.Setenv("CONFIG_STRING_FIELD", "sample")
	os.Setenv("CONFIG_BOOL_FIELD", "true")
	os.Setenv("CONFIG_INT_NESTED_FIELD", "3")
	os.Setenv("CONFIG_STRING_NESTED_FIELD", "nested-sample")

	var config Config
	err = c.Unmarshal(&config)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, config.IntField)
	assert.Equal(t, float64(2), config.FloatField)
	assert.Equal(t, "sample", config.StringField)
	assert.Equal(t, true, config.BoolField)
	assert.Equal(t, 3, config.NestedConfig.IntNestedField)
	assert.Equal(t, "nested-sample", config.NestedConfig.StringNestedField)

}

func TestDefaultConfigure(t *testing.T) {
	c, err := NewViperConfig(
		config.WithDefaults(map[string]interface{}{
			"DEFAULT_KEY": "DEFAULT_VALUE",
		}),
		config.WithTagName("config"),
		config.WithAutomaticEnv(true),
		config.WithConfigFile("sample.env"),
	)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	assert.NotNil(t, c)

	config.SetDefaultConfigure(c)
	var conf Config

	err = config.Unmarshal(&conf)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 9, conf.IntField)
	assert.Equal(t, float64(9.9), conf.FloatField)
	assert.Equal(t, "sample", conf.StringField)
	assert.Equal(t, true, conf.BoolField)
	assert.Equal(t, 9, conf.NestedConfig.IntNestedField)
	assert.Equal(t, "sample", conf.NestedConfig.StringNestedField)

}

func TestDefaultGetConfigureWithEnvFile(t *testing.T) {
	c, err := NewViperConfig(
		config.WithDefaults(map[string]interface{}{
			"DEFAULT_KEY": "DEFAULT_VALUE",
		}),
		config.WithConfigFile("sample.env"),
	)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	ti, err := time.Parse(time.RFC3339, "2020-08-19T15:24:39-07:00")
	if err != nil {
		fmt.Println(err.Error())
	}

	assert.NotNil(t, c)

	config.SetDefaultConfigure(c)
	assert.Equal(t, "config", config.Get("CONFIG"))
	assert.Equal(t, true, config.GetBool("CONFIG_BOOL"))
	assert.Equal(t, "string", config.GetString("CONFIG_STRING"))
	assert.Equal(t, 5, config.GetInt("CONFIG_NUMBER_INT"))
	assert.Equal(t, int32(-999999), config.GetInt32("CONFIG_NUMBER_INT32"))
	assert.Equal(t, int64(999999), config.GetInt64("CONFIG_NUMBER_INT64"))
	assert.Equal(t, uint(999999), config.GetUint("CONFIG_NUMBER_UINT"))
	assert.Equal(t, uint32(999999), config.GetUint32("CONFIG_NUMBER_UINT32"))
	assert.Equal(t, uint64(999999), config.GetUint64("CONFIG_NUMBER_UINT64"))
	assert.Equal(t, 999999.9, config.GetFloat64("CONFIG_NUMBER_FLOAT64"))
	assert.Equal(t, ti, config.GetTime("CONFIG_TIME"))
	assert.Equal(t, time.Second*5, config.GetDuration("CONFIG_TIME_DURATION"))
	assert.Equal(t, time.Second*5, config.GetDuration("CONFIG_TIME_DURATION"))
	// assert.Equal(t, []int{1, 2, 3}, config.GetIntSlice("CONFIG_INT_SLICE"))
	assert.Equal(t, []string{"str1", "str2", "str3"}, config.GetStringSlice("CONFIG_STRING_SLICE"))
	assert.Equal(t, map[string]interface{}{
		"str1": "str1",
		"str2": "str2",
		"str3": "str3",
	}, config.GetStringMap("CONFIG_STRING_MAP"))

}
