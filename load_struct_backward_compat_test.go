package dotenv

import (
	"os"
	"testing"
)

// TestLoadStruct_BackwardCompatibility ensures that existing code without default tags
// continues to work exactly as before
func TestLoadStruct_BackwardCompatibility(t *testing.T) {
	// Clean up any existing environment variables
	defer func() {
		_ = os.Unsetenv("COMPAT_STRING")
		_ = os.Unsetenv("COMPAT_INT")
		_ = os.Unsetenv("COMPAT_BOOL")
		_ = os.Unsetenv("COMPAT_FLOAT")
	}()

	t.Run("Fields without env vars remain zero values", func(t *testing.T) {
		config := &struct {
			StringField string  `env:"COMPAT_STRING"`
			IntField    int     `env:"COMPAT_INT"`
			BoolField   bool    `env:"COMPAT_BOOL"`
			FloatField  float64 `env:"COMPAT_FLOAT"`
		}{}

		err := LoadStruct(config)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// All fields should have zero values
		if config.StringField != "" {
			t.Errorf("Expected StringField to be empty, got: %s", config.StringField)
		}
		if config.IntField != 0 {
			t.Errorf("Expected IntField to be 0, got: %d", config.IntField)
		}
		if config.BoolField != false {
			t.Errorf("Expected BoolField to be false, got: %t", config.BoolField)
		}
		if config.FloatField != 0.0 {
			t.Errorf("Expected FloatField to be 0.0, got: %f", config.FloatField)
		}
	})

	t.Run("Fields with env vars are set correctly", func(t *testing.T) {
		_ = os.Setenv("COMPAT_STRING", "test_value")
		_ = os.Setenv("COMPAT_INT", "42")
		_ = os.Setenv("COMPAT_BOOL", "true")
		_ = os.Setenv("COMPAT_FLOAT", "3.14")

		config := &struct {
			StringField string  `env:"COMPAT_STRING"`
			IntField    int     `env:"COMPAT_INT"`
			BoolField   bool    `env:"COMPAT_BOOL"`
			FloatField  float64 `env:"COMPAT_FLOAT"`
		}{}

		err := LoadStruct(config)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if config.StringField != "test_value" {
			t.Errorf("Expected StringField to be 'test_value', got: %s", config.StringField)
		}
		if config.IntField != 42 {
			t.Errorf("Expected IntField to be 42, got: %d", config.IntField)
		}
		if config.BoolField != true {
			t.Errorf("Expected BoolField to be true, got: %t", config.BoolField)
		}
		if config.FloatField != 3.14 {
			t.Errorf("Expected FloatField to be 3.14, got: %f", config.FloatField)
		}
	})

	t.Run("Mixed fields with and without default tags", func(t *testing.T) {
		_ = os.Setenv("COMPAT_WITH_ENV", "from_env")
		// Intentionally not setting COMPAT_WITHOUT_ENV and COMPAT_WITH_DEFAULT

		config := &struct {
			FieldWithEnv     string `env:"COMPAT_WITH_ENV"`
			FieldWithoutEnv  string `env:"COMPAT_WITHOUT_ENV"`
			FieldWithDefault string `env:"COMPAT_WITH_DEFAULT" default:"default_val"`
		}{}

		err := LoadStruct(config)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if config.FieldWithEnv != "from_env" {
			t.Errorf("Expected FieldWithEnv to be 'from_env', got: %s", config.FieldWithEnv)
		}
		if config.FieldWithoutEnv != "" {
			t.Errorf("Expected FieldWithoutEnv to be empty, got: %s", config.FieldWithoutEnv)
		}
		if config.FieldWithDefault != "default_val" {
			t.Errorf("Expected FieldWithDefault to be 'default_val', got: %s", config.FieldWithDefault)
		}

		_ = os.Unsetenv("COMPAT_WITH_ENV")
	})
}

// TestLoadStruct_DefaultTagEdgeCases tests additional edge cases
func TestLoadStruct_DefaultTagEdgeCases(t *testing.T) {
	t.Run("Default with special characters", func(t *testing.T) {
		config := &struct {
			SpecialChars string `env:"SPECIAL_CHARS" default:"hello@world#2024!"`
			JsonString   string `env:"JSON_STRING" default:"{\"key\":\"value\"}"`
			PathString   string `env:"PATH_STRING" default:"/usr/local/bin:/usr/bin"`
		}{}

		err := LoadStruct(config)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if config.SpecialChars != "hello@world#2024!" {
			t.Errorf("Expected SpecialChars to be 'hello@world#2024!', got: %s", config.SpecialChars)
		}
		if config.JsonString != `{"key":"value"}` {
			t.Errorf("Expected JsonString to be '{\"key\":\"value\"}', got: %s", config.JsonString)
		}
		if config.PathString != "/usr/local/bin:/usr/bin" {
			t.Errorf("Expected PathString to be '/usr/local/bin:/usr/bin', got: %s", config.PathString)
		}
	})

	t.Run("Default with whitespace", func(t *testing.T) {
		config := &struct {
			WhitespaceString string `env:"WHITESPACE" default:"  spaces around  "`
			TabString        string `env:"TAB_STRING" default:"tab	separated"`
			NewlineString    string `env:"NEWLINE_STRING" default:"line\nbreak"`
		}{}

		err := LoadStruct(config)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if config.WhitespaceString != "  spaces around  " {
			t.Errorf("Expected WhitespaceString to preserve spaces, got: '%s'", config.WhitespaceString)
		}
		if config.TabString != "tab	separated" {
			t.Errorf("Expected TabString to contain tab, got: '%s'", config.TabString)
		}
		if config.NewlineString != "line\nbreak" {
			t.Errorf("Expected NewlineString to contain newline, got: '%s'", config.NewlineString)
		}
	})

	t.Run("Default with numeric edge values", func(t *testing.T) {
		config := &struct {
			Zero      int     `env:"ZERO" default:"0"`
			Negative  int     `env:"NEGATIVE" default:"-1"`
			LargeInt  int64   `env:"LARGE_INT" default:"9223372036854775807"`
			SmallInt  int8    `env:"SMALL_INT" default:"-128"`
			TinyFloat float32 `env:"TINY_FLOAT" default:"0.0000001"`
		}{}

		err := LoadStruct(config)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if config.Zero != 0 {
			t.Errorf("Expected Zero to be 0, got: %d", config.Zero)
		}
		if config.Negative != -1 {
			t.Errorf("Expected Negative to be -1, got: %d", config.Negative)
		}
		if config.LargeInt != 9223372036854775807 {
			t.Errorf("Expected LargeInt to be max int64, got: %d", config.LargeInt)
		}
		if config.SmallInt != -128 {
			t.Errorf("Expected SmallInt to be -128, got: %d", config.SmallInt)
		}
		if config.TinyFloat != 0.0000001 {
			t.Errorf("Expected TinyFloat to be 0.0000001, got: %f", config.TinyFloat)
		}
	})

	t.Run("Default with bool edge values", func(t *testing.T) {
		config := &struct {
			BoolTrue   bool `env:"BOOL_TRUE" default:"true"`
			BoolFalse  bool `env:"BOOL_FALSE" default:"false"`
			BoolOne    bool `env:"BOOL_ONE" default:"1"`
			BoolZero   bool `env:"BOOL_ZERO" default:"0"`
			BoolT      bool `env:"BOOL_T" default:"t"`
			BoolF      bool `env:"BOOL_F" default:"f"`
			BoolTUpper bool `env:"BOOL_T_UPPER" default:"T"`
			BoolFUpper bool `env:"BOOL_F_UPPER" default:"F"`
		}{}

		err := LoadStruct(config)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if config.BoolTrue != true {
			t.Errorf("Expected BoolTrue to be true")
		}
		if config.BoolFalse != false {
			t.Errorf("Expected BoolFalse to be false")
		}
		if config.BoolOne != true {
			t.Errorf("Expected BoolOne to be true")
		}
		if config.BoolZero != false {
			t.Errorf("Expected BoolZero to be false")
		}
		if config.BoolT != true {
			t.Errorf("Expected BoolT to be true")
		}
		if config.BoolF != false {
			t.Errorf("Expected BoolF to be false")
		}
		if config.BoolTUpper != true {
			t.Errorf("Expected BoolTUpper to be true")
		}
		if config.BoolFUpper != false {
			t.Errorf("Expected BoolFUpper to be false")
		}
	})
}

// TestLoadStruct_DefaultPrecedence tests the precedence of environment variables over defaults
func TestLoadStruct_DefaultPrecedence(t *testing.T) {
	// Set environment variable
	_ = os.Setenv("PRECEDENCE_TEST", "env_wins")
	defer os.Unsetenv("PRECEDENCE_TEST")

	config := &struct {
		Value string `env:"PRECEDENCE_TEST" default:"default_loses"`
	}{}

	err := LoadStruct(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.Value != "env_wins" {
		t.Errorf("Expected environment variable to override default, got: %s", config.Value)
	}
}

// TestLoadStruct_FieldsWithoutEnvTag tests that fields without env tag are ignored
func TestLoadStruct_FieldsWithoutEnvTag(t *testing.T) {
	config := &struct {
		FieldWithEnv    string `env:"WITH_ENV" default:"has_default"`
		FieldWithoutEnv string `default:"should_be_ignored"`
		FieldPlain      string
	}{}

	err := LoadStruct(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.FieldWithEnv != "has_default" {
		t.Errorf("Expected FieldWithEnv to use default, got: %s", config.FieldWithEnv)
	}
	if config.FieldWithoutEnv != "" {
		t.Errorf("Expected FieldWithoutEnv to remain empty, got: %s", config.FieldWithoutEnv)
	}
	if config.FieldPlain != "" {
		t.Errorf("Expected FieldPlain to remain empty, got: %s", config.FieldPlain)
	}
}

// TestLoadStruct_UnexportedFields tests that unexported fields are safely ignored
func TestLoadStruct_UnexportedFields(t *testing.T) {
	_ = os.Setenv("UNEXPORTED_TEST", "should_not_be_set")
	defer os.Unsetenv("UNEXPORTED_TEST")

	config := &struct {
		ExportedField   string `env:"EXPORTED_TEST" default:"exported_default"`
		unexportedField string `env:"UNEXPORTED_TEST" default:"unexported_default"`
	}{}

	err := LoadStruct(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.ExportedField != "exported_default" {
		t.Errorf("Expected ExportedField to use default, got: %s", config.ExportedField)
	}
	if config.unexportedField != "" {
		t.Errorf("Expected unexportedField to remain empty (unexported), got: %s", config.unexportedField)
	}
}