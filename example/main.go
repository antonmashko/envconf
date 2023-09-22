package main

import (
	"fmt"

	"github.com/antonmashko/envconf"
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
	// envconf.SetLogger(log.New(os.Stdout, "envconf", log.Ldate|log.Ltime))
	if err := envconf.Parse(&cfg); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", cfg)
}
