package dotenv

import (
	"os"
	"testing"
)

func TestParseExportPrefix(t *testing.T) {
	// Clean up env vars before and after test
	cleanup := func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("WITH_SPACES")
		os.Unsetenv("EXTRA_SPACES")
		os.Unsetenv("AFTER_COMMENT")
	}
	cleanup()
	t.Cleanup(cleanup)

	if err := Parse("test/test_export.env"); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "export with simple value",
			key:      "DB_HOST",
			expected: "localhost",
		},
		{
			name:     "export with numeric value",
			key:      "DB_PORT",
			expected: "5432",
		},
		{
			name:     "without export prefix",
			key:      "DB_NAME",
			expected: "mydb",
		},
		{
			name:     "export with double quoted value",
			key:      "DB_USER",
			expected: "admin",
		},
		{
			name:     "export with single quoted value",
			key:      "DB_PASS",
			expected: "secret123",
		},
		{
			name:     "export with spaces in value",
			key:      "WITH_SPACES",
			expected: "value with spaces",
		},
		{
			name:     "export with extra spaces after keyword",
			key:      "EXTRA_SPACES",
			expected: "trimmed",
		},
		{
			name:     "export after comment line",
			key:      "AFTER_COMMENT",
			expected: "works",
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

func TestStripExportPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with export prefix",
			input:    "export KEY=value",
			expected: "KEY=value",
		},
		{
			name:     "with export and extra spaces",
			input:    "export   KEY=value",
			expected: "KEY=value",
		},
		{
			name:     "without export prefix",
			input:    "KEY=value",
			expected: "KEY=value",
		},
		{
			name:     "export in value should not be stripped",
			input:    "KEY=export_data",
			expected: "KEY=export_data",
		},
		{
			name:     "EXPORT uppercase should not be stripped",
			input:    "EXPORT=value",
			expected: "EXPORT=value",
		},
		{
			name:     "export without space is a key",
			input:    "exportKEY=value",
			expected: "exportKEY=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripExportPrefix(tt.input)
			if got != tt.expected {
				t.Errorf("stripExportPrefix(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
