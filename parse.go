package dotenv

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(?m)(\$\{.*?})`)

// Parse parses the .env file located at the given location and set the environment variables.
func Parse(location string) error {
	file, err := os.Open(location)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	items := strings.Split(string(b), linebreak())

	for nb, item := range items {
		if len(item) == 0 || string(item[0]) == "#" {
			continue
		}

		// Remove inline comments
		fields := strings.Fields(item)

		for i, field := range fields {
			if string(field[0]) == "#" {
				fields = fields[:i]
				break
			}
		}

		item = strings.Join(fields, " ")

		split := strings.SplitN(item, "=", 2)
		if len(split) != 2 {
			fmt.Println(split)
			return fmt.Errorf("line %d: cannot get key and value", nb)
		}

		for _, v := range re.FindAllString(split[1], -1) {
			name := strings.TrimRight(strings.TrimLeft(v, "${"), "}")
			split[1] = strings.ReplaceAll(split[1], v, os.Getenv(name))
		}

		if err = os.Setenv(split[0], split[1]); err != nil {
			return err
		}
	}

	return nil
}
