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
	Itr   fmt.Stringer
	Inner *TestInnerStruct `prefix:"INNER_"`
}

func TestStructSerialization_Ok(t *testing.T) {
	os.Setenv("DATA", "foo")
	input := &testStruct{}
	v := depointerize(reflect.ValueOf(input))
	str := newStructFieldData(v, nil, reflect.StructField{})
	err := str.Serialize("")
	if err != nil {
		t.Fatal(err)
	}

	if input.Inner.Data != "foo" {
		t.Fatal("no config")
	}
}
