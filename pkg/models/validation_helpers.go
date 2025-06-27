package models

import "fmt"

// Common validation constants
const (
	MaxKeyLength   = 100
	MaxValueLength = 500
)

// validateKeyValueMap validates a map of key-value pairs
func validateKeyValueMap(kvMap map[string]string, fieldType string) []string {
	var errors []string

	for key, value := range kvMap {
		if key == "" {
			errors = append(errors, fmt.Sprintf("%s keys cannot be empty", fieldType))
		}
		if len(key) > MaxKeyLength {
			errors = append(errors, fmt.Sprintf("%s keys must be 100 characters or less", fieldType))
		}
		if len(value) > MaxValueLength {
			errors = append(errors, fmt.Sprintf("%s values must be 500 characters or less", fieldType))
		}
	}

	return errors
}
