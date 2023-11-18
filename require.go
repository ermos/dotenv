package dotenv

import (
	"fmt"
	"os"
	"strings"
)

// Require check if the given keys are set in the environment variables.
func Require(keys ...string) error {
	var required []string

	for _, key := range keys {
		if os.Getenv(key) == "" {
			required = append(required, key)
		}
	}

	if len(required) > 0 {
		return fmt.Errorf(
			"the following environment variables are required: %s",
			strings.Join(required, ", "),
		)
	}

	return nil
}
