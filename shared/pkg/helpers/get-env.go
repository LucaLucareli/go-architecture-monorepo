package helpers

import (
	"os"
	"strconv"
)

func GetEnv[T any](key string, fallback T) T {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	switch any(fallback).(type) {
	case int:
		if n, err := strconv.Atoi(v); err == nil {
			return any(n).(T)
		}
	case string:
		return any(v).(T)
	}

	return fallback
}
