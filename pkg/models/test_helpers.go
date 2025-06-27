package models

import "testing"

// testSlugValidation is a helper function for testing slug validation functions
func testSlugValidation(t *testing.T, funcName string, validatorFunc func(string) bool, validExamples, invalidExamples []string) {
	t.Helper()

	// Test valid examples
	for _, example := range validExamples {
		t.Run("valid_"+example, func(t *testing.T) {
			if !validatorFunc(example) {
				t.Errorf("%s(%q) = false, want true", funcName, example)
			}
		})
	}

	// Test invalid examples
	for _, example := range invalidExamples {
		t.Run("invalid_"+example, func(t *testing.T) {
			if validatorFunc(example) {
				t.Errorf("%s(%q) = true, want false", funcName, example)
			}
		})
	}
}

// testIsActiveValidation is a helper function for testing IsActive methods
func testIsActiveValidation(t *testing.T, activeGetter, archivedGetter func() bool) {
	t.Helper()

	t.Run("active", func(t *testing.T) {
		if !activeGetter() {
			t.Error("IsActive() = false, want true for non-archived entity")
		}
	})

	t.Run("archived", func(t *testing.T) {
		if archivedGetter() {
			t.Error("IsActive() = true, want false for archived entity")
		}
	})
}

// testStringValidation is a helper function for testing string validation functions
func testStringValidation(t *testing.T, funcName string, validatorFunc func(string) bool, validStrings, invalidStrings []string) {
	t.Helper()

	testSlugValidation(t, funcName, validatorFunc, validStrings, invalidStrings)
}
