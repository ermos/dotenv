package dotenv

import (
	"os"
	"testing"
)

func TestVariableSubstitutionSimpleSyntax(t *testing.T) {
	// Set up environment variables for substitution
	os.Setenv("VERSION", "v2")
	os.Setenv("ENV", "production")
	os.Setenv("PREFIX", "test")
	os.Setenv("VAR1", "hello")
	os.Setenv("VAR2", "world")
	os.Setenv("MY_VAR_NAME", "custom_value")

	// Clean up test vars
	cleanup := func() {
		os.Unsetenv("BASE_URL")
		os.Unsetenv("API_ENDPOINT")
		os.Unsetenv("FULL_PATH")
		os.Unsetenv("LOG_PREFIX")
		os.Unsetenv("START_VAR")
		os.Unsetenv("MULTI")
		os.Unsetenv("WITH_UNDERSCORE")
		os.Unsetenv("DOLLAR_SPACE")
		os.Unsetenv("DOLLAR_NUMBER")
	}
	cleanup()
	t.Cleanup(cleanup)

	if err := Parse("test/test_substitution_simple.env"); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "base value without substitution",
			key:      "BASE_URL",
			expected: "https://api.example.com",
		},
		{
			name:     "simple $VAR syntax",
			key:      "API_ENDPOINT",
			expected: "https://api.example.com/v1",
		},
		{
			name:     "mixed ${VAR} and $VAR syntax",
			key:      "FULL_PATH",
			expected: "https://api.example.com/users/v2",
		},
		{
			name:     "$VAR at end of string",
			key:      "LOG_PREFIX",
			expected: "app_production",
		},
		{
			name:     "${VAR} with delimiter for concatenation",
			key:      "START_VAR",
			expected: "test_value",
		},
		{
			name:     "multiple ${VAR} substitutions with delimiters",
			key:      "MULTI",
			expected: "hello_world",
		},
		{
			name:     "$VAR with underscores in name",
			key:      "WITH_UNDERSCORE",
			expected: "custom_value",
		},
		{
			name:     "dollar followed by space is not substituted",
			key:      "DOLLAR_SPACE",
			expected: "$ not_a_var",
		},
		{
			name:     "dollar followed by number is not substituted",
			key:      "DOLLAR_NUMBER",
			expected: "$123invalid",
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

func TestSubstitutionRegexPatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "braced syntax",
			input:    "prefix_${VAR}_suffix",
			envVars:  map[string]string{"VAR": "value"},
			expected: "prefix_value_suffix",
		},
		{
			name:     "simple syntax",
			input:    "prefix_$VAR_suffix",
			envVars:  map[string]string{"VAR_suffix": "replaced"},
			expected: "prefix_replaced",
		},
		{
			name:     "simple syntax end of string",
			input:    "prefix_$VAR",
			envVars:  map[string]string{"VAR": "value"},
			expected: "prefix_value",
		},
		{
			name:     "both syntaxes",
			input:    "$START/${MIDDLE}/$END",
			envVars:  map[string]string{"START": "a", "MIDDLE": "b", "END": "c"},
			expected: "a/b/c",
		},
		{
			name:     "undefined variable braced",
			input:    "value_${UNDEFINED}",
			envVars:  map[string]string{},
			expected: "value_",
		},
		{
			name:     "undefined variable simple",
			input:    "value_$UNDEFINED",
			envVars:  map[string]string{},
			expected: "value_",
		},
		{
			name:     "dollar not followed by identifier",
			input:    "price $5",
			envVars:  map[string]string{},
			expected: "price $5",
		},
		{
			name:     "underscore start",
			input:    "$_PRIVATE",
			envVars:  map[string]string{"_PRIVATE": "secret"},
			expected: "secret",
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
