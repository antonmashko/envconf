package envconf

import (
	"os"
	"reflect"
	"testing"
)

type TestInnerStruct struct {
	Data string `env:"DATA"`
}

type testStruct struct {
	Inner ***TestInnerStruct `prefix:"INNER_"`
}

func TestStructSerialization_Ok(t *testing.T) {
	os.Setenv("DATA", "foo")
	input := &testStruct{}
	// v := depointerize(reflect.ValueOf(input))
	cfg := &structType{}
	err := cfg.Init(reflect.ValueOf(input).Elem(), nil, reflect.StructField{})
	if err != nil {
		t.Fatal(err)
	}

	err = cfg.Define()
	if err != nil {
		t.Fatal(err)
	}

	if (**input.Inner).Data != "foo" {
		t.Fatal("no config")
	}
}
