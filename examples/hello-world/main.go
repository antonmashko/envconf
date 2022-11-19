package main

import (
	"fmt"

	"github.com/antonmashko/envconf"
)

type cfg struct {
	Outer *struct {
		Inner struct {
			Data string `flag:"test" env:"TEST" default:"test" required:"1"`
		}
	}
}

func main() {
	var c cfg
	err := envconf.Parse(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.Outer.Inner.Data)
}
