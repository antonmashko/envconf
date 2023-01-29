package main

import (
	"fmt"

	"github.com/antonmashko/envconf"
)

type Example struct {
	Inner struct {
		FlagField   string `flag:"*"`
		EnvVarField string `env:"*"`
	}
}

// Run `go run main.go --help` for getting help output with auto-generated names
func main() {
	var cfg Example
	if err := envconf.Parse(&cfg); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", cfg)
}
