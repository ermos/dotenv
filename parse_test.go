package dotenv

import (
	"os"
	"testing"
)

func TestParseOK(t *testing.T) {
	if err := Parse("test/.env"); err != nil {
		t.Error(err)
	}
	if os.Getenv("toto") != "tata" {
		t.Error("toto is not equal to tata")
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
