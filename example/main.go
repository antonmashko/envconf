package main

import (
	"fmt"
	"log"

	"github.com/antonmashko/envconf"
	"github.com/antonmashko/envconf/option"
)

type Example struct {
	Field1 string `flag:"flag-name" env:"ENV_VAR_NAME" default:"default-value"`
	Inner  struct {
		FlagField   string `flag:"*"`
		EnvVarField string `env:"*"`
	}
	Password string `default:"pass1"`
	API      struct {
		Token string `default:"token1"`
	}
}

// Run `go run main.go --help` for getting help output with auto-generated names
func main() {
	var cfg Example
	if err := envconf.Parse(&cfg, option.WithLog(log.Default())); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", cfg)
}
