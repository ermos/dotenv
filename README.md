# DotEnv

> Loads environment variables from .env file with Go

## Installation

```bash
go get github.com/ermos/dotenv
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/ermos/dotenv"
)

func main() {
    if err := dotenv.Parse(".env"); err != nil {
        log.Fatal(err)
    }

    fmt.Println(os.Getenv("DATABASE_URL"))
}
```

## .env File Syntax

This parser follows the [dotenv](https://github.com/bkeepers/dotenv) conventions and supports most features from popular implementations ([Node.js dotenv](https://github.com/motdotla/dotenv), [Python python-dotenv](https://github.com/theskumar/python-dotenv)).

### Basic Syntax

```bash
# Simple key-value
KEY=value

# Comments (lines starting with #)
# This is a comment
KEY=value  # Inline comments are supported
```

### Quoted Values

```bash
# Double quotes (escape sequences supported)
MESSAGE="Hello\nWorld"

# Single quotes (literal, no escape processing)
PATH='/usr/local/bin'

# Quotes preserve spaces
GREETING="Hello, World!"
```

### Export Prefix

Compatible with shell scripts - the `export` keyword is stripped:

```bash
export DATABASE_URL=postgres://localhost/mydb
export API_KEY="secret"
```

> Reference: [POSIX Shell Command Language](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html)

### Variable Substitution

```bash
# Braced syntax (recommended)
BASE_URL=https://api.example.com
API_ENDPOINT=${BASE_URL}/v1/users

# Simple syntax (POSIX-style)
HOME_DIR=/home/user
CONFIG=$HOME_DIR/.config

# Note: $VAR_NAME matches the longest valid identifier
# Use ${VAR}_suffix to delimit variable names
PREFIX=${APP}_production
```

> Reference: [POSIX Parameter Expansion](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_06_02)

### Default Values

```bash
# Use default if variable is unset OR empty
DATABASE_HOST=${DB_HOST:-localhost}
DATABASE_PORT=${DB_PORT:-5432}

# Use default only if variable is unset (keeps empty values)
OPTIONAL=${MAYBE_EMPTY-default}

# Nested defaults
CONFIG_PATH=${CUSTOM_PATH:-${HOME:-/tmp}/config}
```

| Syntax | Behavior |
|--------|----------|
| `${VAR:-default}` | Use `default` if `VAR` is **unset or empty** |
| `${VAR-default}` | Use `default` only if `VAR` is **unset** |

> Reference: [Bash Parameter Expansion](https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html)

### Multiline Values

```bash
# Quoted multiline (quotes span multiple lines)
PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
-----END RSA PRIVATE KEY-----"

# JSON configuration
CONFIG='{
  "host": "localhost",
  "port": 8080
}'

# Backslash line continuation
LONG_COMMAND=docker run \
  --name myapp \
  --env-file .env \
  myimage:latest
```

### Escape Sequences (in double quotes)

| Sequence | Result |
|----------|--------|
| `\n` | Newline |
| `\t` | Tab |
| `\r` | Carriage return |
| `\\` | Backslash |
| `\"` | Double quote |

```bash
MESSAGE="Line 1\nLine 2\tTabbed"
PATH="C:\\Users\\name"
QUOTED="He said \"Hello\""
```

### Syntax Summary

| Feature | Example | Supported |
|---------|---------|-----------|
| Basic assignment | `KEY=value` | Yes |
| Comments | `# comment` | Yes |
| Inline comments | `KEY=value # comment` | Yes |
| Double quotes | `KEY="value"` | Yes |
| Single quotes | `KEY='value'` | Yes |
| Export prefix | `export KEY=value` | Yes |
| Variable substitution | `${VAR}` or `$VAR` | Yes |
| Default values | `${VAR:-default}` | Yes |
| Multiline (quoted) | `KEY="line1\nline2"` | Yes |
| Backslash continuation | `KEY=a\`<br>`b` | Yes |

## API Reference

### Parse

Parses a `.env` file and sets environment variables.

```go
err := dotenv.Parse(".env")
err := dotenv.Parse("/path/to/custom.env")
```

### Require

Validates that required environment variables are set.

```go
err := dotenv.Require("DATABASE_URL", "API_KEY", "SECRET")
if err != nil {
    log.Fatal(err) // "the following environment variables are required: API_KEY, SECRET"
}
```

### LoadStruct

Loads environment variables into a struct using tags.

```go
type Config struct {
    Host     string `env:"HOST" default:"localhost"`
    Port     int    `env:"PORT" default:"8080"`
    Debug    bool   `env:"DEBUG"`
    APIKey   string `env:"API_KEY" required:"true"`
}

var cfg Config
if err := dotenv.LoadStruct(&cfg); err != nil {
    log.Fatal(err)
}
```

#### Struct Tags

| Tag | Description |
|-----|-------------|
| `env:"VAR_NAME"` | Maps field to environment variable |
| `default:"value"` | Default value if not set |
| `required:"true"` | Error if variable is not set |
| `validator:"name"` | Custom validator (with `LoadStructWithOptions`) |

#### Supported Types

`string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`

#### Custom Validators

```go
opts := dotenv.LoadOptions{
    Validators: map[string]dotenv.Validator{
        "email": func(v reflect.Value) error {
            if !strings.Contains(v.String(), "@") {
                return fmt.Errorf("invalid email format")
            }
            return nil
        },
    },
}

type Config struct {
    Email string `env:"EMAIL" validator:"email"`
}

var cfg Config
err := dotenv.LoadStructWithOptions(&cfg, opts)
```

### Typed Getters

Helper functions to retrieve and convert environment variables.

```go
// Basic getters (return zero value on error)
dotenv.GetString("KEY")
dotenv.GetInt("PORT")
dotenv.GetBool("DEBUG")
dotenv.GetFloat64("RATE")
dotenv.GetFloat32("RATIO")
dotenv.GetInt64("BIG_NUMBER")
dotenv.GetUint("COUNT")

// With default values
dotenv.GetStringOrDefault("KEY", "default")
dotenv.GetIntOrDefault("PORT", 8080)
dotenv.GetBoolOrDefault("DEBUG", false)
dotenv.GetFloat64OrDefault("RATE", 1.0)
dotenv.GetFloat32OrDefault("RATIO", 0.5)
dotenv.GetInt64OrDefault("BIG_NUMBER", 0)
dotenv.GetUintOrDefault("COUNT", 1)
```

## License

MIT
