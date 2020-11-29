# DotEnv ♻️
> Loads environment variables from .env file with Go

## Installation

```bash
go get github.com/ermos/dotenv
```

## Usage

This package contains only one method : ``dotenv.Parse()``,
It allows to load a list of environment variable from an `.env` file.

## Example

#### .env
```.env
toto=tata
lib_desc=Loads environment variables from .env file with Go
```

#### main.go
```go
import (
 "fmt"
 "github.com/ermos/dotenv"
 "os"
)

func main () {
    // Load environment variables
    if err := dotenv.Parse(".env"); err != nil {
        log.Fatal(err)
    }   
    
    fmt.Println(os.Getenv("lib_desc"))
    fmt.Printf("Toto is equal to %s", os.Getenv("toto"))
}
```

#### Result
```bash
> Loads environment variables from .env file with Go
> Toto is equal to tata
```
