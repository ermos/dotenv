package dotenv

import "testing"

func TestRequireOK(t *testing.T) {
	if err := Parse("test/.env"); err != nil {
		t.Error(err)
	}

	err := Require(
		"NAME",
		"DESCRIPTION",
	)
	if err != nil {
		t.Error(err)
	}
}

func TestRequireNOK(t *testing.T) {
	if err := Parse("test/.env"); err != nil {
		t.Error(err)
	}

	err := Require(
		"NAME",
		"DESCRIPTION",
		"INVALID",
	)
	if err == nil {
		t.Errorf("INVALID is not define but return nil")
	}
}
