// Package autoenv provides a simple and flexible way to load environment variables into Go structs.
//
// It supports automatic mapping of environment variables to struct fields using reflection,
// with support for both JSON and ENV tags, nested structs, and various data types including slices.
// The package is designed to have zero dependencies and focuses on being both powerful and easy to use.
//
// Basic Usage:
//
//	type Config struct {
//	    Port    int      `json:"port"`    // Will look for PORT
//	    Host    string   `json:"host"`    // Will look for HOST
//	    Debug   bool     `json:"debug"`   // Will look for DEBUG
//	}
//
//	cfg := &Config{}
//	autoenv.Load(cfg)
//
// Tag Support:
//
// The package supports both json and env tags:
//
//	type Config struct {
//	    // Using JSON tags (automatically converted to SNAKE_CASE)
//	    DatabaseURL string `json:"databaseUrl"` // Will look for DATABASE_URL
//
//	    // Using ENV tags (exact match)
//	    APIKey     string `env:"API_KEY"`      // Will look for API_KEY
//	}
//
// Features:
//
//   - Automatic SNAKE_CASE conversion for field names
//   - Support for nested structs
//   - Optional .env file loading
//   - Environment variable prefixing
//   - Field ignoring capabilities
//   - Custom logging support
//   - Slice support (comma-separated values)
//
// Supported Types:
//   - string
//   - bool
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - []string, []bool, []int, []uint, []float64
//
// Configuration Options:
//
//	autoenv.Load(cfg,
//	    autoenv.WithPrefix("APP"),        // Add prefix to all env vars
//	    autoenv.WithEnvFile(),            // Load from .env file
//	    autoenv.WithVerbose(),            // Enable verbose logging
//	    autoenv.WithOnlyEnvTag(),         // Only use env tags
//	    autoenv.WithIgnore("debug"),      // Ignore specific fields
//	)
//
// For more information and examples, visit: https://go.g3deon.com/autoenv
package autoenv
