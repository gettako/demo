package services

import "strconv"

// MapConfig is a simple in-memory implementation of contracts.Config.
// It is intended to be registered as a Singleton in the IoC container.
type MapConfig struct {
	data map[string]any
}

// NewMapConfig creates a new configuration populated with initial values.
func NewMapConfig(initial map[string]any) *MapConfig {
	return &MapConfig{
		data: initial,
	}
}

func (c *MapConfig) String(key string) string {
	if val, ok := c.data[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
		// Fallback to basic conversion
		return ""
	}
	return ""
}

func (c *MapConfig) Int(key string) int {
	if val, ok := c.data[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

func (c *MapConfig) Bool(key string) bool {
	if val, ok := c.data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
