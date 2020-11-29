package dotenv

import (
	"fmt"
	"os"
)

// Allow to require an environment variable which is necessary for your golang application
func Require(keys ...string) error {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			return fmt.Errorf("%s's environment variable is required", key)
		}
	}
	return nil
}