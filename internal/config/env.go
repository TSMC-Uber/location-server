package config

import (
	"os"
)

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func MustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic("env variable " + key + " is not set")
	}
	return value
}
