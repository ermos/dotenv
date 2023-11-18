package dotenv

import (
	"os"
	"testing"
)

func TestParseOK(t *testing.T) {
	if err := Parse("test/.env"); err != nil {
		t.Error(err)
	}
	if os.Getenv("VERSION") != "1.0.0" {
		t.Errorf("VERSION should be 1.0.0")
	}

	if os.Getenv("NAME") != "dotenv" {
		t.Errorf("NAME should be dotenv")
	}

	if os.Getenv("DESCRIPTION") != "Loads environment variables from .env file with Go" {
		t.Errorf("DESCRIPTION should be Loads environment variables from .env file with Go, but is %s", os.Getenv("DESCRIPTION"))
	}
}

func TestParseOpen(t *testing.T) {
	if err := Parse("test/not-exist/.env"); err == nil {
		t.Errorf("file doesnt exist but return nil")
	}
}

func TestParseNotEqual(t *testing.T) {
	if err := Parse("test/not-equal.env"); err == nil {
		t.Errorf(".env line 2 doesnt contains equal sign but OK ?")
	}
}

func TestParseNoKey(t *testing.T) {
	if err := Parse("test/no-key.env"); err == nil {
		t.Errorf(".env line 2 doesnt contains key but OK ?")
	}
}
