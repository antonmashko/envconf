package envconf_test

import (
	"testing"

	"github.com/antonmashko/envconf"
)

type stringTextUnmarshaller string

func (tu *stringTextUnmarshaller) UnmarshalText(text []byte) error {
	*tu = stringTextUnmarshaller("txt")
	return nil
}

func TestParse_StringTextUnmarshaller_Ok(t *testing.T) {
	tc := struct {
		Field1 stringTextUnmarshaller `default:"10"`
	}{}
	err := envconf.Parse(&tc)
	if err != nil {
		t.Fatal(err)
	}
	if tc.Field1 != "txt" {
		t.Fatal("unexpected result: ", tc.Field1)
	}
}

func TestParse_PointerStringTextUnmarshaller_Ok(t *testing.T) {
	tc := struct {
		Field1 *stringTextUnmarshaller `default:"10"`
	}{}
	err := envconf.Parse(&tc)
	if err != nil {
		t.Fatal(err)
	}
	if *tc.Field1 != "txt" {
		t.Fatal("unexpected result: ", *tc.Field1)
	}
}

func TestParse_DoublePointerStringTextUnmarshaller_Ok(t *testing.T) {
	tc := struct {
		Field1 **stringTextUnmarshaller `default:"10"`
	}{}
	err := envconf.Parse(&tc)
	if err != nil {
		t.Fatal(err)
	}
	if **tc.Field1 != "txt" {
		t.Fatal("unexpected result: ", **tc.Field1)
	}
}

type structTextUnmarshaller struct {
	data    string
	invoked bool
}

func (tu *structTextUnmarshaller) UnmarshalText(text []byte) error {
	tu.invoked = true
	tu.data = "txt"
	return nil
}

func TestParse_StructTextUnmarshaller_Ok(t *testing.T) {
	tc := struct {
		Field1 structTextUnmarshaller `default:"10"`
	}{}
	err := envconf.Parse(&tc)
	if err != nil {
		t.Fatal(err)
	}
	if tc.Field1.data != "txt" || !tc.Field1.invoked {
		t.Fatal("unexpected result: ", tc.Field1)
	}
}

func TestParse_PointerStructTextUnmarshaller_Ok(t *testing.T) {
	tc := struct {
		Field1 *structTextUnmarshaller `default:"10"`
	}{}
	err := envconf.Parse(&tc)
	if err != nil {
		t.Fatal(err)
	}
	if tc.Field1.data != "txt" || !tc.Field1.invoked {
		t.Fatal("unexpected result: ", tc.Field1)
	}
}

func TestParse_DoublePointerStructTextUnmarshaller_Ok(t *testing.T) {
	tc := struct {
		Field1 **structTextUnmarshaller `default:"10"`
	}{}
	err := envconf.Parse(&tc)
	if err != nil {
		t.Fatal(err)
	}
	if (*tc.Field1).data != "txt" || !(*tc.Field1).invoked {
		t.Fatal("unexpected result: ", *tc.Field1)
	}
}
