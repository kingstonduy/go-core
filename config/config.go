package config

import "time"

const (
	EnvFileKey        = "GO_ENV_FILE"
	DefaultConfigFile = ".env"
)

var (
	DefaultConfigure = newNoopsConfigure()
)

func SetDefaultConfigure(config Configure) {
	DefaultConfigure = config
}

func Get(key string) interface{} {
	return DefaultConfigure.Get(key)
}

func GetString(key string) string {
	return DefaultConfigure.GetString(key)
}

func GetBool(key string) bool {
	return DefaultConfigure.GetBool(key)
}

func GetInt(key string) int {
	return DefaultConfigure.GetInt(key)
}

func GetInt32(key string) int32 {
	return DefaultConfigure.GetInt32(key)
}

func GetInt64(key string) int64 {
	return DefaultConfigure.GetInt64(key)
}

func GetUint(key string) uint {
	return DefaultConfigure.GetUint(key)
}

func GetUint16(key string) uint16 {
	return DefaultConfigure.GetUint16(key)
}

func GetUint32(key string) uint32 {
	return DefaultConfigure.GetUint32(key)
}

func GetUint64(key string) uint64 {
	return DefaultConfigure.GetUint64(key)
}

func GetFloat64(key string) float64 {
	return DefaultConfigure.GetFloat64(key)
}

func GetTime(key string) time.Time {
	return DefaultConfigure.GetTime(key)
}

func GetDuration(key string) time.Duration {
	return DefaultConfigure.GetDuration(key)
}

func GetIntSlice(key string) []int {
	return DefaultConfigure.GetIntSlice(key)
}

func GetStringSlice(key string) []string {
	return DefaultConfigure.GetStringSlice(key)
}

func GetStringMap(key string) map[string]any {
	return DefaultConfigure.GetStringMap(key)
}

func Unmarshal(dest interface{}) error {
	return DefaultConfigure.Unmarshal(dest)
}

func UnmarshalKey(key string, dest any) error {
	return DefaultConfigure.UnmarshalKey(key, dest)
}

type Configure interface {
	Get(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint16(key string) uint16
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]any
	Unmarshal(dest interface{}) error
	UnmarshalKey(key string, dest any) error
}
