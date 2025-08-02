package main

import (
	"fmt"
	"go.g3deon.com/autoenv"
	"os"
)

type AppConfig struct {
	Port int
	Host string
	DB   struct {
		Host string
		Port int
		User string
		Pass string
		Name string
	}
}

func main() {
	var config AppConfig

	loader := autoenv.NewLoader(
		autoenv.WithPrefix("MY_APP"),
		autoenv.WithPaths([]string{".env.example"}),
		autoenv.WithVerbose(),
	)
	if err := loader.Load(&config); err != nil {
		fmt.Printf("failed to load config: %s", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", config)
}
