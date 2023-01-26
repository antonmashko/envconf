package envconf_test

import (
	"testing"

	"github.com/antonmashko/envconf"
)

func TestParseFieldInNestedStruct_Ok(t *testing.T) {
	data := struct {
		Inner struct {
			Field1 int `default:"123"`
		}
	}{}
	err := envconf.Parse(&data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Inner.Field1 != 123 {
		t.Fatalf("incorrect value. expected=123 actual=%d", data.Inner.Field1)
	}
}

func TestParseFieldInNNestedStruct_Ok(t *testing.T) {
	data := struct {
		Inner1 struct {
			Inner12 struct {
			}
			Inner13 struct {
				Inner131 struct {
					Field1 int `default:"123"`
				}
			}
		}
		Inner2 struct{}
	}{}
	err := envconf.Parse(&data)
	if err != nil {
		t.Fatal(err)
	}
	if data.Inner1.Inner13.Inner131.Field1 != 123 {
		t.Fatalf("incorrect value. expected=123 actual=%d", data.Inner1.Inner13.Inner131.Field1)
	}
}
