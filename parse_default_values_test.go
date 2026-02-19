package dotenv

import (
	"os"
	"testing"
)

func TestDefaultValueSubstitution(t *testing.T) {
	// Set up environment
	os.Setenv("DEFINED_VAR", "defined_value")
	os.Setenv("EMPTY_VAR", "")

	// Ensure undefined vars are unset
	os.Unsetenv("UNDEFINED_HOST")
	os.Unsetenv("UNDEFINED_PORT")
	os.Unsetenv("UNSET_VAR")
	os.Unsetenv("LEVEL1")
	os.Unsetenv("LEVEL2")
	os.Unsetenv("UNDEFINED")
	os.Unsetenv("UNDEFINED_URL")
	os.Unsetenv("UNDEFINED_EMPTY")
	os.Unsetenv("UNDEFINED_QUOTED")

	// Clean up test vars
	cleanup := func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("EXISTING_VAR")
		os.Unsetenv("EMPTY_CHECK")
		os.Unsetenv("UNSET_ONLY")
		os.Unsetenv("EMPTY_UNSET_ONLY")
		os.Unsetenv("NESTED")
		os.Unsetenv("SPECIAL_DEFAULT")
		os.Unsetenv("URL_DEFAULT")
		os.Unsetenv("EMPTY_DEFAULT")
		os.Unsetenv("QUOTED_DEFAULT")
		os.Unsetenv("DEFINED_VAR")
		os.Unsetenv("EMPTY_VAR")
	}
	cleanup()

	// Re-set after cleanup
	os.Setenv("DEFINED_VAR", "defined_value")
	os.Setenv("EMPTY_VAR", "")

	t.Cleanup(cleanup)

	if err := Parse("test/test_default_values.env"); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "default value for undefined variable",
			key:      "DB_HOST",
			expected: "localhost",
		},
		{
			name:     "default numeric value for undefined variable",
			key:      "DB_PORT",
			expected: "5432",
		},
		{
			name:     "existing variable ignores default",
			key:      "EXISTING_VAR",
			expected: "defined_value",
		},
		{
			name:     "empty variable uses default with :-",
			key:      "EMPTY_CHECK",
			expected: "default_for_empty",
		},
		{
			name:     "unset variable uses default with -",
			key:      "UNSET_ONLY",
			expected: "default_unset",
		},
		{
			name:     "empty variable keeps empty with - (no colon)",
			key:      "EMPTY_UNSET_ONLY",
			expected: "",
		},
		{
			name:     "nested defaults",
			key:      "NESTED",
			expected: "final_default",
		},
		{
			name:     "default with spaces",
			key:      "SPECIAL_DEFAULT",
			expected: "hello world",
		},
		{
			name:     "default with URL",
			key:      "URL_DEFAULT",
			expected: "https://example.com/path?q=1",
		},
		{
			name:     "empty default value",
			key:      "EMPTY_DEFAULT",
			expected: "",
		},
		{
			name:     "default in quoted value",
			key:      "QUOTED_DEFAULT",
			expected: "quoted default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := os.Getenv(tt.key)
			if got != tt.expected {
				t.Errorf("os.Getenv(%q) = %q, want %q", tt.key, got, tt.expected)
			}
		})
	}
}

func TestExtractDefaultValue(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		envVars         map[string]string
		expected        string
		expectedVarName string
	}{
		{
			name:     "simple default with :-",
			input:    "${VAR:-default}",
			envVars:  map[string]string{},
			expected: "default",
		},
		{
			name:     "simple default with -",
			input:    "${VAR-default}",
			envVars:  map[string]string{},
			expected: "default",
		},
		{
			name:     "existing var ignores default with :-",
			input:    "${VAR:-default}",
			envVars:  map[string]string{"VAR": "value"},
			expected: "value",
		},
		{
			name:     "empty var uses default with :-",
			input:    "${VAR:-default}",
			envVars:  map[string]string{"VAR": ""},
			expected: "default",
		},
		{
			name:     "empty var keeps empty with -",
			input:    "${VAR-default}",
			envVars:  map[string]string{"VAR": ""},
			expected: "",
		},
		{
			name:     "no default syntax",
			input:    "${VAR}",
			envVars:  map[string]string{},
			expected: "",
		},
		{
			name:     "default with special chars",
			input:    "${URL:-https://example.com}",
			envVars:  map[string]string{},
			expected: "https://example.com",
		},
		{
			name:     "default with spaces",
			input:    "${MSG:-hello world}",
			envVars:  map[string]string{},
			expected: "hello world",
		},
		{
			name:     "empty default",
			input:    "${VAR:-}",
			envVars:  map[string]string{},
			expected: "",
		},
		{
			name:     "default containing hyphen",
			input:    "${VAR:-my-value}",
			envVars:  map[string]string{},
			expected: "my-value",
		},
		{
			name:     "default containing colon",
			input:    "${VAR:-host:port}",
			envVars:  map[string]string{},
			expected: "host:port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			t.Cleanup(func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			})

			got := processSubstitution(tt.input)
			if got != tt.expected {
				t.Errorf("processSubstitution(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
