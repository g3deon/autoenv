# AutoEnv for Go

[![Actions Status](https://github.com/g3deon/autoenv/actions/workflows/go.yml/badge.svg)](https://github.com/g3deon/autoenv/actions)
[![Go Reference](https://pkg.go.dev/badge/go.g3deon.com/autoenv.svg)](https://pkg.go.dev/go.g3deon.com/autoenv)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Automatically map environment variables into Go structs using reflection, tag-based mapping, and optional .env support. Built for modern Go services that require fast and declarative configuration loading — with zero dependencies.

## Installation

```sh
go get go.g3deon.com/autoenv
```

## Usage

```go
import "go.g3deon.com/autoenv"

type Config struct {
  Port  int    `json:"port"`
  Debug bool   `json:"debug"`
  Host  string `json:"host"`
}

func main() {
  cfg := &Config{}
  autoenv.Load(cfg)
}
```

Using custom `env` tag:

```go
type Config struct {
  Port int    `env:"PORT"`
  Debug bool  `json:"debug"`
}
```

> [!IMPORTANT]
> When using the `env` tag, no SNAKE\_CASE conversion is applied. The key is taken exactly as written.

## Embedded structs

```go
type DBConfig struct {
  Host string `json:"dbHost"`
  Port int    `json:"dbPort"`
}

type Config struct {
  DBConfig
  AppName string `json:"appName"`
}
```

## .env file support

Load environment variables from a `.env` file:

```go
autoenv.Load(cfg, autoenv.WithEnvFile())
```

Custom path:

```go
autoenv.Load(cfg, autoenv.WithEnvFilePath("config/.env"))
```

Supported format:

```sh
# Standard KEY=VALUE syntax
PORT=3000                     
# Bash-style export statement
export DEBUG=true             
# Double-quoted string
DATABASE_URL="postgres://localhost/db"  
# Single-quoted string
API_KEY='key'                 
```

## Ignore fields

Unexported fields are ignored by default. To ignore specific exported fields:

```go
autoenv.Load(cfg, autoenv.WithIgnore("debug"))
```

Combine multiple:

```go
autoenv.Load(
  cfg,
  autoenv.WithIgnore("field"),
  autoenv.WithIgnore("other_field"),
)
```

Replace the full ignore list:

```go
autoenv.Load(
  cfg,
  autoenv.WithIgnores([]string{"field1", "field2"}),
)
```

## Only use `env` tags

```go
autoenv.Load(cfg, autoenv.WithOnlyEnvTag())
```

## Verbose mode

```go
autoenv.Load(cfg, autoenv.WithVerbose())
```

Example output:

```sh
[AUTOENV] Reflecting field: Port
[AUTOENV] Mapped to env key: PORT
[AUTOENV] Value not found in environment
```

## Custom Logger

```go
type MyLogger struct{}

func (l *MyLogger) Debugf(format string, args ...any)   {}
func (l *MyLogger) Errorf(format string, args ...any)   {}

autoenv.Load(cfg, autoenv.WithLogger(&MyLogger{}))
```


## License

MIT © 2025 G3deon, Inc.
