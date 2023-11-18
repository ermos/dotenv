package dotenv

import (
	"os"
	"testing"
)

func TestGetString(t *testing.T) {
	err := os.Setenv("TEST", "test")
	if err != nil {
		t.Error(err)
	}

	if GetString("TEST") != "test" {
		t.Error("GetString() should return \"test\"")
	}
}

func TestGetEmptyString(t *testing.T) {
	err := os.Setenv("TEST", "")
	if err != nil {
		t.Error(err)
	}

	if GetString("TEST") != "" {
		t.Error("GetString() should return \"\"")
	}
}

func TestGetStringWithDefault(t *testing.T) {
	err := os.Setenv("TEST", "")
	if err != nil {
		t.Error(err)
	}

	if GetStringOrDefault("TEST", "default") != "default" {
		t.Error("GetString() should return \"default\"")
	}
}

func TestGetStringWithDefaultWithValidValue(t *testing.T) {
	err := os.Setenv("TEST", "test")
	if err != nil {
		t.Error(err)
	}

	if GetStringOrDefault("TEST", "default") != "test" {
		t.Error("GetString() should return \"test\"")
	}
}

func TestGetInt(t *testing.T) {
	err := os.Setenv("TEST", "1")
	if err != nil {
		t.Error(err)
	}

	if GetInt("TEST") != 1 {
		t.Error("GetInt() should return 1")
	}
}

func TestGetIntWithInvalidValue(t *testing.T) {
	err := os.Setenv("TEST", "invalid")
	if err != nil {
		t.Error(err)
	}

	if GetInt("TEST") != 0 {
		t.Error("GetInt() should return 0")
	}
}

func TestGetIntWithDefault(t *testing.T) {
	err := os.Setenv("TEST", "")
	if err != nil {
		t.Error(err)
	}

	if GetIntOrDefault("TEST", 1) != 1 {
		t.Error("GetInt() should return 1")
	}
}

func TestGetIntWithDefaultWithValidValue(t *testing.T) {
	err := os.Setenv("TEST", "1")
	if err != nil {
		t.Error(err)
	}

	if GetIntOrDefault("TEST", 2) != 1 {
		t.Error("GetInt() should return 1")
	}
}

func TestGetBool(t *testing.T) {
	err := os.Setenv("TEST", "true")
	if err != nil {
		t.Error(err)
	}

	if GetBool("TEST") != true {
		t.Error("GetBool() should return true")
	}
}

func TestGetBoolWithInvalidValue(t *testing.T) {
	err := os.Setenv("TEST", "invalid")
	if err != nil {
		t.Error(err)
	}

	if GetBool("TEST") != false {
		t.Error("GetBool() should return false")
	}
}

func TestGetBoolWithDefault(t *testing.T) {
	err := os.Setenv("TEST", "")
	if err != nil {
		t.Error(err)
	}

	if GetBoolOrDefault("TEST", true) != true {
		t.Error("GetBool() should return true")
	}
}

func TestGetBoolWithDefaultWithValidValue(t *testing.T) {
	err := os.Setenv("TEST", "true")
	if err != nil {
		t.Error(err)
	}

	if GetBoolOrDefault("TEST", false) != true {
		t.Error("GetBool() should return true")
	}
}

func TestGetFloat64(t *testing.T) {
	err := os.Setenv("TEST", "1.1")
	if err != nil {
		t.Error(err)
	}

	if GetFloat64("TEST") != 1.1 {
		t.Error("GetFloat64() should return 1.1")
	}
}

func TestGetFloat64WithInvalidValue(t *testing.T) {
	err := os.Setenv("TEST", "invalid")
	if err != nil {
		t.Error(err)
	}

	if GetFloat64("TEST") != 0 {
		t.Error("GetFloat64() should return 0")
	}
}

func TestGetFloat64WithDefault(t *testing.T) {
	err := os.Setenv("TEST", "")
	if err != nil {
		t.Error(err)
	}

	if GetFloat64OrDefault("TEST", 1.1) != 1.1 {
		t.Error("GetFloat64() should return 1.1")
	}
}

func TestGetFloat64WithDefaultWithValidValue(t *testing.T) {
	err := os.Setenv("TEST", "1.1")
	if err != nil {
		t.Error(err)
	}

	if GetFloat64OrDefault("TEST", 2.2) != 1.1 {
		t.Error("GetFloat64() should return 1.1")
	}
}
