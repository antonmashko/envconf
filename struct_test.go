package envconf

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type TestInnerStruct struct {
	Data string `env:"DATA"`
}

type testStruct struct {
	Str   fmt.Stringer
	Inner ***TestInnerStruct `prefix:"INNER_"`
}

func TestStructSerialization_Ok(t *testing.T) {
	os.Setenv("DATA", "foo")
	input := &testStruct{}
	// v := depointerize(reflect.ValueOf(input))
	cfg := newStructType(reflect.ValueOf(input).Elem(), nil, reflect.StructField{})
	err := cfg.Init()
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
