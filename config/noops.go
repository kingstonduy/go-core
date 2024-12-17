package config

import (
	"log"
	"time"
)

type noopsConfigure struct{}

func newNoopsConfigure() Configure {
	return &noopsConfigure{}
}

// Get implements Configure.
func (n *noopsConfigure) Get(key string) interface{} {
	n.noopsWarning()
	return nil
}

// GetBool implements Configure.
func (n *noopsConfigure) GetBool(key string) bool {
	n.noopsWarning()
	return false
}

// GetDuration implements Configure.
func (n *noopsConfigure) GetDuration(key string) time.Duration {
	n.noopsWarning()
	return 0
}

// GetFloat64 implements Configure.
func (n *noopsConfigure) GetFloat64(key string) float64 {
	n.noopsWarning()
	return 0
}

// GetInt implements Configure.
func (n *noopsConfigure) GetInt(key string) int {
	n.noopsWarning()
	return 0
}

// GetInt32 implements Configure.
func (n *noopsConfigure) GetInt32(key string) int32 {
	n.noopsWarning()
	return 0
}

// GetInt64 implements Configure.
func (n *noopsConfigure) GetInt64(key string) int64 {
	n.noopsWarning()
	return 0
}

// GetIntSlice implements Configure.
func (n *noopsConfigure) GetIntSlice(key string) []int {
	n.noopsWarning()
	return nil
}

// GetString implements Configure.
func (n *noopsConfigure) GetString(key string) string {
	n.noopsWarning()
	return ""
}

// GetStringMap implements Configure.
func (n *noopsConfigure) GetStringMap(key string) map[string]any {
	n.noopsWarning()
	return nil
}

// GetStringSlice implements Configure.
func (n *noopsConfigure) GetStringSlice(key string) []string {
	n.noopsWarning()
	return nil
}

// GetTime implements Configure.
func (n *noopsConfigure) GetTime(key string) time.Time {
	n.noopsWarning()
	return time.Time{}
}

// GetUint implements Configure.
func (n *noopsConfigure) GetUint(key string) uint {
	n.noopsWarning()
	return 0
}

// GetUint16 implements Configure.
func (n *noopsConfigure) GetUint16(key string) uint16 {
	n.noopsWarning()
	return 0
}

// GetUint32 implements Configure.
func (n *noopsConfigure) GetUint32(key string) uint32 {
	n.noopsWarning()
	return 0
}

// GetUint64 implements Configure.
func (n *noopsConfigure) GetUint64(key string) uint64 {
	n.noopsWarning()
	return 0
}

// Unmarshal implements Configure.
func (n *noopsConfigure) Unmarshal(dest interface{}) error {
	n.noopsWarning()
	return nil
}

// UnmarshalKey implements Configure.
func (n *noopsConfigure) UnmarshalKey(key string, dest any) error {
	n.noopsWarning()
	return nil
}

func (n *noopsConfigure) noopsWarning() {
	log.Print("[WARN] No default configure was set. Using noops configure as default. Set the default configure to do all functions\n")
}
