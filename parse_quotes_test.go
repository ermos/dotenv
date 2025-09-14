package dotenv

import (
	"os"
	"testing"
)

func TestParseQuotedValues(t *testing.T) {
	// Clear environment first
	os.Clearenv()

	if err := Parse("test/test_quotes.env"); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	tests := []struct {
		key      string
		expected string
		desc     string
	}{
		{"SIMPLE", "value", "simple value"},
		{"WITH_SPACES", "value with spaces", "unquoted value with spaces"},
		{"QUOTED_SPACES", "value with spaces", "double quoted value with spaces"},
		{"SINGLE_QUOTED", "another value with spaces", "single quoted value with spaces"},
		{"MIXED_QUOTES", "it's a test", "value with apostrophe in double quotes"},
		{"URL", "https://example.com/path", "URL with inline comment"},
		{"PATH_WITH_SPACES", "/home/user/my documents/file.txt", "path with spaces"},
		{"QUOTED_PATH", "/home/user/my documents/file.txt", "quoted path with spaces"},
		{"EMPTY_VALUE", "", "empty value"},
		{"QUOTED_EMPTY", "", "quoted empty value"},
		{"TRAILING_SPACE", "value", "value with trailing space should be trimmed"},
		{"LEADING_SPACE", "value", "value with leading space should be trimmed"},
		{"BOTH_SPACES", "value", "value with both spaces should be trimmed"},
		{"INLINE_COMMENT", "value#comment", "value with immediate comment (# without space is part of value)"},
		{"SPACED_COMMENT", "value", "value with spaced comment"},
		{"ESCAPED_QUOTE", `value with " quote`, "escaped quote in quoted value"},
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
