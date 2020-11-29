package dotenv

import "testing"

func TestRequireOK(t *testing.T) {
	if err := Parse("test/.env"); err != nil {
		t.Error(err)
	}

	err := Require(
		"toto",
		"lib_desc",
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
		"toto",
		"lib_desc",
		"tata",
	)
	if err == nil {
		t.Errorf("tata is not define but available")
	}
}
