# AutoEnv for Go

[![Actions Status](https://github.com/g3deon/autoenv/actions/workflows/go.yml/badge.svg)](https://github.com/g3deon/autoenv/actions)
[![Go Reference](https://pkg.go.dev/badge/go.g3deon.com/autoenv.svg)](https://pkg.go.dev/go.g3deon.com/autoenv)
[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Automatically map environment variables into Go structs using reflection and tag-based mapping, with optional .env file
support. Designed for modern Go services that need fast, declarative, zero-dependency configuration loading. It includes
automatic SNAKE_CASE conversion, support for JSON and env tags, embedded structs, custom logging, prefix handling, field
ignoring, and even parsing of slices and common data types—all without external libraries.

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

## Tag Support

### JSON Tags

By default, the library uses JSON tags and converts them to SNAKE_CASE:

```go
type Config struct {
	DatabaseURL string `json:"databaseUrl"` // Will look for DATABASE_URL 
	MaxRetries int `json:"maxRetries"`      // Will look for MAX_RETRIES
}
```

### ENV Tags

Use `env` tags for custom environment variable names:

```go
type Config struct {
	Port int `env:"SERVICE_PORT"` // Will look for SERVICE_PORT exactly 
	Debug bool `env:"APP_DEBUG"`  // Will look for APP_DEBUG exactly 
}
```

## Advanced Features

### Prefix Support

Add a prefix to all environment variables:

```go 
autoenv.Load(cfg, autoenv.WithPrefix("APP"))
```

With the above configuration:

- `Port` will look for `APP_PORT`
- `Debug` will look for `APP_DEBUG`

### Slice Support

The library supports slices of basic types:

```go
type Config struct {
	Ports []int `json:"ports"` // PORTS=8080,8081,8082 
	Hosts []string `json:"hosts"` // HOSTS=localhost,127.0.0.1 
	Flags []bool `json:"flags"`   // FLAGS=true,false,true 
	Numbers []float64 `json:"numbers"` // NUMBERS=1.1,2.2,3.3 
}
```

### Nested Structs

```go
type DatabaseConfig struct {
	Host string `json:"host"` // DATABASE_HOST
	Port int `json:"port"`    // DATABASE_PORT
	Password string `json:"password"` // DATABASE_PASSWORD
}

type Config struct {
	Database DatabaseConfig
	AppName string `json:"appName"` // APP_NAME
}
```

### Custom Logger Interface

```go 
type Logger interface {
	DebugF(format string, args ...any)
	ErrorF(format string, args ...any)
	// ...
}

// Implement your own logger type MyLogger struct{}
func (l *MyLogger) DebugF(format string, args ...any) {
	log.Printf("[DEBUG] "+format, args...)
}

func (l *MyLogger) ErrorF(format string, args ...any) {
	log.Printf("[ERROR] "+format, args...)
}

// Use your logger 
autoenv.Load(cfg, autoenv.WithLogger(&MyLogger{}))
```

## License

MIT © 2025 G3deon, Inc.
