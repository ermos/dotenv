package dotenv

import (
	"os"
	"testing"
)

func TestVariableSubstitution(t *testing.T) {
	// Set up base environment variables
	_ = os.Setenv("DOMAIN", "example.com")
	_ = os.Setenv("YEAR", "2024")

	// Clear other vars to test
	_ = os.Unsetenv("BASE_PATH")
	_ = os.Unsetenv("LOG_PATH")
	_ = os.Unsetenv("DATA_PATH")
	_ = os.Unsetenv("FULL_URL")
	_ = os.Unsetenv("CONFIG_PATH")
	_ = os.Unsetenv("NESTED")
	_ = os.Unsetenv("EMPTY_SUB")
	_ = os.Unsetenv("PARTIAL")
	_ = os.Unsetenv("UNDEFINED_VAR")

	if err := Parse("test/test_substitution.env"); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	tests := []struct {
		key      string
		expected string
		desc     string
	}{
		{"BASE_PATH", "/var/app", "base path without substitution"},
		{"LOG_PATH", "/var/app/logs", "simple substitution"},
		{"DATA_PATH", "/var/app/data", "another simple substitution"},
		{"FULL_URL", "https://example.com/api", "substitution from pre-existing env var"},
		{"CONFIG_PATH", "/var/app/config", "substitution in quoted value"},
		{"NESTED", "/var/app/logs/archive/2024", "nested substitution"},
		{"EMPTY_SUB", "", "substitution of undefined variable"},
		{"PARTIAL", "prefix_/var/app_suffix", "substitution in middle of value"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			actual := os.Getenv(tt.key)
			if actual != tt.expected {
				t.Errorf("For %s: expected %q, got %q", tt.key, tt.expected, actual)
			}
		})
	}
}
