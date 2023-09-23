package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/antonmashko/envconf"
	"github.com/antonmashko/envconf/option"
)

type Example struct {
	Field1 string `flag:"flag-name" env:"ENV_VAR_NAME" default:"default-value"`
	Field2 string `json:"field-2"`
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
	jsonConf := &envconf.Json{}
	configFilePath := flag.String("config", "./conf.json", "")
	fpOpt := option.WithFlagParsed(func() error {
		b, err := os.ReadFile(*configFilePath)
		if err != nil {
			return err
		}
		*jsonConf = envconf.Json(b)
		return nil
	})
	err := envconf.ParseWithExternal(&cfg, jsonConf, option.WithLog(log.Default()), fpOpt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", cfg)
}
