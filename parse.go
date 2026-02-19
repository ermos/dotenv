package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// Compile regex once at package level for better performance
var re = regexp.MustCompile(`(?m)(\$\{.*?})`)

// Parse parses the .env file located at the given location and set the environment variables.
// This function now properly handles:
// - Quoted values (both single and double quotes)
// - Values with spaces
// - Inline comments with proper detection
// - Escape sequences in quoted strings
// - Leading/trailing whitespace trimming
func Parse(location string) error {
	file, err := os.Open(location)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines and full-line comments
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		// Strip "export " prefix if present
		line = stripExportPrefix(line)

		// Find the first = sign
		equalIndex := strings.Index(line, "=")
		if equalIndex == -1 {
			return fmt.Errorf("line %d: cannot get key and value", lineNum)
		}

		// Extract key and value
		key := strings.TrimSpace(line[:equalIndex])
		if key == "" {
			return fmt.Errorf("line %d: cannot get key and value", lineNum)
		}

		// Get raw value (everything after =)
		rawValue := line[equalIndex+1:]

		// Process the value
		value := processEnvValue(rawValue)

		// Variable substitution
		for _, v := range re.FindAllString(value, -1) {
			name := strings.TrimRight(strings.TrimLeft(v, "${"), "}")
			value = strings.ReplaceAll(value, v, os.Getenv(name))
		}

		if err = os.Setenv(key, value); err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

// processEnvValue handles quote removal, inline comments, and trimming
func processEnvValue(raw string) string {
	// Trim leading whitespace to check for quotes
	trimmed := strings.TrimLeftFunc(raw, unicode.IsSpace)

	// Handle quoted values
	if len(trimmed) > 0 {
		switch trimmed[0] {
		case '"':
			return extractQuoted(trimmed, '"')
		case '\'':
			return extractQuoted(trimmed, '\'')
		}
	}

	// For unquoted values, handle inline comments
	value := removeInlineComments(raw)

	// Trim spaces from unquoted values
	return strings.TrimSpace(value)
}

// extractQuoted extracts value from within quotes, handling escape sequences
func extractQuoted(s string, quote byte) string {
	if len(s) < 2 {
		return s
	}

	var result strings.Builder
	escaped := false

	for i := 1; i < len(s); i++ {
		ch := s[i]

		if escaped {
			// Handle common escape sequences
			switch ch {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case quote:
				result.WriteByte(quote)
			default:
				// Keep the backslash for unknown escapes
				result.WriteByte('\\')
				result.WriteByte(ch)
			}
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			continue
		}

		if ch == quote {
			// Found closing quote - return result
			return result.String()
		}

		result.WriteByte(ch)
	}

	// Unclosed quote - return what we have
	return result.String()
}

// removeInlineComments removes inline comments from unquoted values
func removeInlineComments(s string) string {
	// Look for # that's either at start or preceded by whitespace
	for i := 0; i < len(s); i++ {
		if s[i] == '#' {
			if i == 0 || unicode.IsSpace(rune(s[i-1])) {
				return s[:i]
			}
		}
	}
	return s
}

// stripExportPrefix removes the "export " prefix from a line if present.
// Only strips lowercase "export" followed by at least one space.
func stripExportPrefix(line string) string {
	const prefix = "export"
	trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)

	if len(trimmed) > len(prefix) && trimmed[:len(prefix)] == prefix && unicode.IsSpace(rune(trimmed[len(prefix)])) {
		return strings.TrimLeftFunc(trimmed[len(prefix):], unicode.IsSpace)
	}

	return line
}
