package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Parse parses the .env file located at the given location and set the environment variables.
// This function now properly handles:
// - Quoted values (both single and double quotes)
// - Multiline values within quotes
// - Backslash line continuation
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

		// Handle multiline values and backslash continuation
		rawValue, lineNum, err = readFullValue(rawValue, scanner, &lineNum)
		if err != nil {
			return err
		}

		// Process the value
		value := processEnvValue(rawValue)

		// Variable substitution
		value = processSubstitution(value)

		if err = os.Setenv(key, value); err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

// readFullValue reads a complete value that may span multiple lines.
// Handles quoted multiline values and backslash line continuation.
func readFullValue(rawValue string, scanner *bufio.Scanner, lineNum *int) (string, int, error) {
	trimmed := strings.TrimLeftFunc(rawValue, unicode.IsSpace)

	// Check for quoted multiline value
	if len(trimmed) > 0 && (trimmed[0] == '"' || trimmed[0] == '\'') {
		quote := trimmed[0]
		// Check if the closing quote is on the same line
		if !hasClosingQuote(trimmed[1:], quote) {
			// Need to read more lines
			var builder strings.Builder
			builder.WriteString(rawValue)

			for scanner.Scan() {
				*lineNum++
				nextLine := scanner.Text()
				builder.WriteByte('\n')
				builder.WriteString(nextLine)

				if hasClosingQuote(nextLine, quote) {
					break
				}
			}
			return builder.String(), *lineNum, nil
		}
		return rawValue, *lineNum, nil
	}

	// Check for backslash line continuation (unquoted values)
	if strings.HasSuffix(strings.TrimRightFunc(rawValue, unicode.IsSpace), "\\") {
		var builder strings.Builder
		currentValue := rawValue

		for {
			trimmedVal := strings.TrimRightFunc(currentValue, unicode.IsSpace)
			if !strings.HasSuffix(trimmedVal, "\\") {
				builder.WriteString(currentValue)
				break
			}
			// Remove trailing backslash
			builder.WriteString(trimmedVal[:len(trimmedVal)-1])

			if !scanner.Scan() {
				break
			}
			*lineNum++
			currentValue = scanner.Text()
		}
		return builder.String(), *lineNum, nil
	}

	return rawValue, *lineNum, nil
}

// hasClosingQuote checks if a string contains an unescaped closing quote.
func hasClosingQuote(s string, quote byte) bool {
	escaped := false
	for i := 0; i < len(s); i++ {
		if escaped {
			escaped = false
			continue
		}
		if s[i] == '\\' {
			escaped = true
			continue
		}
		if s[i] == quote {
			return true
		}
	}
	return false
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

// processSubstitution replaces variable references with their values.
// Supports both ${VAR} and $VAR syntax, with optional default values.
// Default value syntax:
//   - ${VAR:-default} : use default if VAR is unset or empty
//   - ${VAR-default}  : use default only if VAR is unset
func processSubstitution(value string) string {
	var result strings.Builder
	i := 0

	for i < len(value) {
		if value[i] == '$' && i+1 < len(value) {
			if value[i+1] == '{' {
				// ${VAR} syntax - find matching closing brace
				content, end := extractBracedContent(value, i+2)
				if end > i {
					result.WriteString(resolveWithDefault(content))
					i = end
					continue
				}
			} else if isIdentifierStart(value[i+1]) {
				// $VAR syntax
				end := i + 2
				for end < len(value) && isIdentifierChar(value[end]) {
					end++
				}
				name := value[i+1 : end]
				result.WriteString(os.Getenv(name))
				i = end
				continue
			}
		}
		result.WriteByte(value[i])
		i++
	}

	return result.String()
}

// extractBracedContent extracts content from ${...}, handling nested braces.
// Returns the content and the index after the closing brace.
func extractBracedContent(s string, start int) (string, int) {
	depth := 1
	i := start

	for i < len(s) && depth > 0 {
		if s[i] == '{' {
			depth++
		} else if s[i] == '}' {
			depth--
		}
		if depth > 0 {
			i++
		}
	}

	if depth == 0 {
		return s[start:i], i + 1
	}
	return "", start
}

// isIdentifierStart checks if a byte can start a variable name.
func isIdentifierStart(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '_'
}

// isIdentifierChar checks if a byte can be part of a variable name.
func isIdentifierChar(b byte) bool {
	return isIdentifierStart(b) || (b >= '0' && b <= '9')
}

// resolveWithDefault handles ${VAR}, ${VAR:-default}, and ${VAR-default} syntax.
func resolveWithDefault(content string) string {
	// Check for :- (use default if unset OR empty)
	if idx := strings.Index(content, ":-"); idx != -1 {
		name := content[:idx]
		defaultVal := content[idx+2:]
		val, exists := os.LookupEnv(name)
		if !exists || val == "" {
			return processSubstitution(defaultVal) // Allow nested substitution
		}
		return val
	}

	// Check for - (use default only if unset)
	if idx := strings.Index(content, "-"); idx != -1 {
		name := content[:idx]
		defaultVal := content[idx+1:]
		val, exists := os.LookupEnv(name)
		if !exists {
			return processSubstitution(defaultVal) // Allow nested substitution
		}
		return val
	}

	// No default syntax, just get the variable
	return os.Getenv(content)
}
