package main

import "github.com/antonmashko/envconf"

type example struct {
	Value1 int    `flag:"value" default:"10"`
	Value2 string `flag:"-" env:"ENV_VALUE" description:"some description"`
	Qwerty Qwerty
}

type Qwerty struct {
	Value3 float32 `env:"QWERTY_VALUE" required:"true"`
}

func main() {
	var e example
	envconf.ParseStruct(&e)
}
