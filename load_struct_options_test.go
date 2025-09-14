package dotenv

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoadStructWithOptions_FormatValidation(t *testing.T) {
	// Set up test environment variables
	_ = os.Setenv("TEST_EMAIL_A", "user@example.com")
	_ = os.Setenv("TEST_EMAIL_B", "Arnold@example.com")

	defer func() {
		_ = os.Unsetenv("TEST_EMAIL_A")
		_ = os.Unsetenv("TEST_EMAIL_B")
	}()

	configA := &struct {
		Email string `env:"TEST_EMAIL_A" validator:"startWithA"`
	}{}

	configB := &struct {
		Email string `env:"TEST_EMAIL_B" validator:"startWithA"`
	}{}

	opts := LoadOptions{
		Validators: map[string]Validator{
			"startWithA": func(value reflect.Value) error {
				if value.Kind() != reflect.String {
					return fmt.Errorf("startWithA validator can only be used with string fields")
				}

				if !strings.HasPrefix(value.String(), "A") {
					return fmt.Errorf("value must start with 'A'")
				}

				return nil
			},
		},
	}

	err := LoadStructWithOptions(configA, opts)
	if err == nil {
		t.Fatalf("Expected error for invalid email, got nil")
	}

	// Verify values are loaded correctly
	if configA.Email != "user@example.com" {
		t.Errorf("Expected email user@example.com, got %s", configA.Email)
	}

	err = LoadStructWithOptions(configB, opts)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify values are loaded correctly
	if configB.Email != "Arnold@example.com" {
		t.Errorf("Expected email Arnold@example.com, got %s", configB.Email)
	}
}
