package dotenv

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoadStruct_DefaultTag(t *testing.T) {
	// Clean up any existing environment variables
	defer func() {
		_ = os.Unsetenv("TEST_STRING")
		_ = os.Unsetenv("TEST_INT")
		_ = os.Unsetenv("TEST_BOOL")
	}()

	config := &struct {
		StringWithDefault    string  `env:"TEST_STRING" default:"default_value"`
		StringWithoutDefault string  `env:"TEST_MISSING"`
		IntWithDefault       int     `env:"TEST_INT" default:"42"`
		BoolWithDefault      bool    `env:"TEST_BOOL" default:"true"`
		FloatWithDefault     float64 `env:"TEST_FLOAT" default:"3.14"`
	}{}

	err := LoadStruct(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check default values are applied
	if config.StringWithDefault != "default_value" {
		t.Errorf("Expected StringWithDefault to be 'default_value', got: %s", config.StringWithDefault)
	}

	if config.StringWithoutDefault != "" {
		t.Errorf("Expected StringWithoutDefault to be empty, got: %s", config.StringWithoutDefault)
	}

	if config.IntWithDefault != 42 {
		t.Errorf("Expected IntWithDefault to be 42, got: %d", config.IntWithDefault)
	}

	if config.BoolWithDefault != true {
		t.Errorf("Expected BoolWithDefault to be true, got: %t", config.BoolWithDefault)
	}

	if config.FloatWithDefault != 3.14 {
		t.Errorf("Expected FloatWithDefault to be 3.14, got: %f", config.FloatWithDefault)
	}
}

func TestLoadStruct_EnvOverridesDefault(t *testing.T) {
	// Set environment variables that should override defaults
	_ = os.Setenv("TEST_STRING_OVERRIDE", "env_value")
	_ = os.Setenv("TEST_INT_OVERRIDE", "100")
	_ = os.Setenv("TEST_BOOL_OVERRIDE", "false")

	defer func() {
		_ = os.Unsetenv("TEST_STRING_OVERRIDE")
		_ = os.Unsetenv("TEST_INT_OVERRIDE")
		_ = os.Unsetenv("TEST_BOOL_OVERRIDE")
	}()

	config := &struct {
		StringField string `env:"TEST_STRING_OVERRIDE" default:"default_string"`
		IntField    int    `env:"TEST_INT_OVERRIDE" default:"42"`
		BoolField   bool   `env:"TEST_BOOL_OVERRIDE" default:"true"`
	}{}

	err := LoadStruct(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Environment values should override defaults
	if config.StringField != "env_value" {
		t.Errorf("Expected StringField to be 'env_value', got: %s", config.StringField)
	}

	if config.IntField != 100 {
		t.Errorf("Expected IntField to be 100, got: %d", config.IntField)
	}

	if config.BoolField != false {
		t.Errorf("Expected BoolField to be false, got: %t", config.BoolField)
	}
}

func TestLoadStruct_DefaultWithInvalidType(t *testing.T) {
	config := &struct {
		IntField int `env:"TEST_INVALID_INT" default:"not_a_number"`
	}{}

	err := LoadStruct(config)
	if err == nil {
		t.Fatal("Expected error for invalid int default, got nil")
	}

	expectedErrMsg := "failed to parse int field IntField"
	if err.Error()[:len(expectedErrMsg)] != expectedErrMsg {
		t.Errorf("Expected error message to start with '%s', got: %s", expectedErrMsg, err.Error())
	}
}

func TestLoadStruct_DefaultEmptyString(t *testing.T) {
	config := &struct {
		EmptyDefault string `env:"TEST_EMPTY" default:""`
		NoDefault    string `env:"TEST_MISSING"`
	}{}

	err := LoadStruct(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Empty default should still set the field
	if config.EmptyDefault != "" {
		t.Errorf("Expected EmptyDefault to be empty string, got: %s", config.EmptyDefault)
	}

	// No default should leave field untouched
	if config.NoDefault != "" {
		t.Errorf("Expected NoDefault to be empty (zero value), got: %s", config.NoDefault)
	}
}

func TestLoadStruct_NestedStructWithDefaults(t *testing.T) {
	config := &struct {
		Database struct {
			Host string `env:"DB_HOST" default:"localhost"`
			Port int    `env:"DB_PORT" default:"5432"`
		}
		Cache struct {
			TTL int `env:"CACHE_TTL" default:"3600"`
		}
	}{}

	err := LoadStruct(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.Database.Host != "localhost" {
		t.Errorf("Expected Database.Host to be 'localhost', got: %s", config.Database.Host)
	}

	if config.Database.Port != 5432 {
		t.Errorf("Expected Database.Port to be 5432, got: %d", config.Database.Port)
	}

	if config.Cache.TTL != 3600 {
		t.Errorf("Expected Cache.TTL to be 3600, got: %d", config.Cache.TTL)
	}
}

func TestLoadStruct_DefaultWithUnsignedInts(t *testing.T) {
	config := &struct {
		Uint8Field  uint8  `env:"TEST_UINT8" default:"255"`
		Uint16Field uint16 `env:"TEST_UINT16" default:"65535"`
		Uint32Field uint32 `env:"TEST_UINT32" default:"4294967295"`
		Uint64Field uint64 `env:"TEST_UINT64" default:"18446744073709551615"`
	}{}

	err := LoadStruct(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.Uint8Field != 255 {
		t.Errorf("Expected Uint8Field to be 255, got: %d", config.Uint8Field)
	}

	if config.Uint16Field != 65535 {
		t.Errorf("Expected Uint16Field to be 65535, got: %d", config.Uint16Field)
	}

	if config.Uint32Field != 4294967295 {
		t.Errorf("Expected Uint32Field to be 4294967295, got: %d", config.Uint32Field)
	}

	if config.Uint64Field != 18446744073709551615 {
		t.Errorf("Expected Uint64Field to be 18446744073709551615, got: %d", config.Uint64Field)
	}
}

func TestLoadStructWithOptions_DefaultAndValidator(t *testing.T) {
	config := &struct {
		Email string `env:"TEST_EMAIL_DEFAULT" default:"admin@example.com" validator:"hasAt"`
	}{}

	opts := LoadOptions{
		Validators: map[string]Validator{
			"hasAt": func(value reflect.Value) error {
				if value.Kind() != reflect.String {
					return fmt.Errorf("hasAt validator can only be used with string fields")
				}
				if !strings.Contains(value.String(), "@") {
					return fmt.Errorf("value must contain '@' symbol")
				}
				return nil
			},
		},
	}

	err := LoadStructWithOptions(config, opts)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.Email != "admin@example.com" {
		t.Errorf("Expected Email to be 'admin@example.com', got: %s", config.Email)
	}
}

func TestLoadStructWithOptions_DefaultFailsValidator(t *testing.T) {
	config := &struct {
		Email string `env:"TEST_EMAIL_INVALID" default:"invalid-email" validator:"hasAt"`
	}{}

	opts := LoadOptions{
		Validators: map[string]Validator{
			"hasAt": func(value reflect.Value) error {
				if value.Kind() != reflect.String {
					return fmt.Errorf("hasAt validator can only be used with string fields")
				}
				if !strings.Contains(value.String(), "@") {
					return fmt.Errorf("value must contain '@' symbol")
				}
				return nil
			},
		},
	}

	err := LoadStructWithOptions(config, opts)
	if err == nil {
		t.Fatal("Expected validation error for invalid default email, got nil")
	}

	expectedErrMsg := "format validation failed for field Email"
	if err.Error()[:len(expectedErrMsg)] != expectedErrMsg {
		t.Errorf("Expected error message to start with '%s', got: %s", expectedErrMsg, err.Error())
	}
}