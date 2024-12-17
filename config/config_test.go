package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultGetConfigureWithEnvFile(t *testing.T) {
	assert.Equal(t, nil, Get("CONFIG"))
	assert.Equal(t, false, GetBool("CONFIG_BOOL"))
	assert.Equal(t, "", GetString("CONFIG_STRING"))
	assert.Equal(t, 0, GetInt("CONFIG_NUMBER_INT"))
	assert.Equal(t, int32(0), GetInt32("CONFIG_NUMBER_INT32"))
	assert.Equal(t, int64(0), GetInt64("CONFIG_NUMBER_INT64"))
	assert.Equal(t, uint(0), GetUint("CONFIG_NUMBER_UINT"))
	assert.Equal(t, uint32(0), GetUint32("CONFIG_NUMBER_UINT32"))
	assert.Equal(t, uint64(0), GetUint64("CONFIG_NUMBER_UINT64"))
	assert.Equal(t, float64(0), GetFloat64("CONFIG_NUMBER_FLOAT64"))
	assert.Equal(t, time.Time{}, GetTime("CONFIG_TIME"))
	assert.Equal(t, time.Duration(0), GetDuration("CONFIG_TIME_DURATION"))
	assert.Equal(t, time.Duration(0), GetDuration("CONFIG_TIME_DURATION"))
}
