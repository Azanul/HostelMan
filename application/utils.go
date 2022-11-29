package main

import (
	"os"
)

func Getenv(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
