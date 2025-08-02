package main

import (
	"fmt"
	"go.g3deon.com/autoenv"
	"os"
)

type AppConfig struct {
	Port int
	Host string
}

func main() {
	var config AppConfig

	os.Setenv("PORT", "8080")
	os.Setenv("HOST", "localhost")

	if err := autoenv.Load(&config); err != nil {
		fmt.Printf("failed to load config: %s", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", config)
}
