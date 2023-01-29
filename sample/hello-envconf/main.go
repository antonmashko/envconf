package main

import (
	"fmt"

	"github.com/antonmashko/envconf"
)

type Example struct {
	Field1 string `flag:"flag-name" env:"ENV_VAR_NAME" default:"default-value"`
}

func main() {
	var cfg Example
	if err := envconf.Parse(&cfg); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", cfg)
}
