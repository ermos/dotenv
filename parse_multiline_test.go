package dotenv

import (
	"os"
	"testing"
)

func TestMultilineValues(t *testing.T) {
	// Clean up test vars
	cleanup := func() {
		os.Unsetenv("MULTILINE_SIMPLE")
		os.Unsetenv("MULTILINE_SINGLE")
		os.Unsetenv("CERTIFICATE")
		os.Unsetenv("JSON_CONFIG")
		os.Unsetenv("WITH_EMPTY_LINES")
		os.Unsetenv("REGULAR_AFTER")
		os.Unsetenv("BASE_PATH")
		os.Unsetenv("SCRIPT")
		os.Unsetenv("CONTINUED")
	}
	cleanup()
	t.Cleanup(cleanup)

	if err := Parse("test/test_multiline.env"); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "simple multiline with double quotes",
			key:      "MULTILINE_SIMPLE",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "multiline with single quotes",
			key:      "MULTILINE_SINGLE",
			expected: "first\nsecond\nthird",
		},
		{
			name:     "certificate-like content",
			key:      "CERTIFICATE",
			expected: "-----BEGIN CERTIFICATE-----\nMIIBkTCB+wIJAKHBfpeg\n-----END CERTIFICATE-----",
		},
		{
			name:     "JSON content with escaped quotes",
			key:      "JSON_CONFIG",
			expected: "{\n  \"key\": \"value\",\n  \"nested\": {\n    \"foo\": \"bar\"\n  }\n}",
		},
		{
			name:     "multiline with empty lines",
			key:      "WITH_EMPTY_LINES",
			expected: "start\n\nmiddle\n\nend",
		},
		{
			name:     "regular value after multiline",
			key:      "REGULAR_AFTER",
			expected: "normal_value",
		},
		{
			name:     "multiline with variable substitution",
			key:      "SCRIPT",
			expected: "#!/bin/bash\ncd /app\necho done",
		},
		{
			name:     "backslash line continuation",
			key:      "CONTINUED",
			expected: "first_partsecond_partthird_part",
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

func TestExtractMultilineQuoted(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		quote    byte
		expected string
	}{
		{
			name:     "single line",
			input:    "\"simple\"",
			quote:    '"',
			expected: "simple",
		},
		{
			name:     "multiline",
			input:    "\"line1\nline2\"",
			quote:    '"',
			expected: "line1\nline2",
		},
		{
			name:     "multiline with empty line",
			input:    "\"start\n\nend\"",
			quote:    '"',
			expected: "start\n\nend",
		},
		{
			name:     "multiline single quotes",
			input:    "'line1\nline2'",
			quote:    '\'',
			expected: "line1\nline2",
		},
		{
			name:     "unclosed quote multiline",
			input:    "\"line1\nline2",
			quote:    '"',
			expected: "line1\nline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractQuoted(tt.input, tt.quote)
			if got != tt.expected {
				t.Errorf("extractQuoted(%q, %q) = %q, want %q", tt.input, tt.quote, got, tt.expected)
			}
		})
	}
}
