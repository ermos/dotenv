package dotenv

import (
	"os"
	"strconv"
)

// GetString returns the value in string format of the environment variable named by the key.
func GetString(key string) string {
	return os.Getenv(key)
}

// GetStringOrDefault returns the value in string format or the default value if the environment variable is not set.
func GetStringOrDefault(key, defaultValue string) string {
	value := GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetInt returns the value in int format of the environment variable named by the key.
func GetInt(key string) int {
	result, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}
	return result
}

// GetIntOrDefault returns the value in int format or the default value if the environment variable is not set.
func GetIntOrDefault(key string, defaultValue int) int {
	result, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return result
}

// GetBool returns the value in bool format of the environment variable named by the key.
func GetBool(key string) bool {
	result, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false
	}
	return result
}

// GetBoolOrDefault returns the value in bool format or the default value if the environment variable is not set.
func GetBoolOrDefault(key string, defaultValue bool) bool {
	result, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return result
}

// GetFloat64 returns the value in float64 format of the environment variable named by the key.
func GetFloat64(key string) float64 {
	result, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil {
		return 0
	}
	return result
}

// GetFloat64OrDefault returns the value in float64 format or the default value if the environment variable is not set.
func GetFloat64OrDefault(key string, defaultValue float64) float64 {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return result
}

// GetFloat32 returns the value in float32 format of the environment variable named by the key.
func GetFloat32(key string) float32 {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}

	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0
	}

	return float32(result)
}

// GetFloat32OrDefault returns the value in float32 format or the default value.
func GetFloat32OrDefault(key string, defaultValue float32) float32 {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultValue
	}

	return float32(result)
}

// GetUint returns the value in uint format of the environment variable named by the key.
func GetUint(key string) uint {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}

	result, err := strconv.ParseUint(value, 10, 0)
	if err != nil {
		return 0
	}

	return uint(result)
}

// GetUintOrDefault returns the value in uint format or the default value.
func GetUintOrDefault(key string, defaultValue uint) uint {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	result, err := strconv.ParseUint(value, 10, 0)
	if err != nil {
		return defaultValue
	}

	return uint(result)
}

// GetInt64 returns the value in int64 format of the environment variable named by the key.
func GetInt64(key string) int64 {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}

	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}

	return result
}

// GetInt64OrDefault returns the value in int64 format or the default value.
func GetInt64OrDefault(key string, defaultValue int64) int64 {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return result
}
