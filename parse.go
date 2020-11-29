package dotenv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Parse(location string) error {
	file, err := os.Open(location)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	items := strings.Split(string(b),"\n")
	if err = file.Close(); err != nil {
		return err
	}

	for l, item := range items {
		if len(item) != 0 && string(item[0]) != "#" {
			split := strings.Split(item, "=")
			if len(split) != 2 {
				return fmt.Errorf("line %d: cannot get key and value", l)
			}

			err = os.Setenv(split[0], split[1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
