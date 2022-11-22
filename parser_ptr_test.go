package envconf_test

import (
	"testing"

	"github.com/antonmashko/envconf"
)

func TestStringPointer_InitFromDefault_Ok(t *testing.T) {
	data := struct {
		Field *string `default:"test"`
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	if data.Field == nil || *data.Field != "test" {
		t.Fatalf("incorrect value. expected=test actual=%s", *data.Field)
	}
}

func TestStructPointer_WithFieldFromDefault_Ok(t *testing.T) {
	data := struct {
		Inner *struct {
			Field string `default:"test"`
		}
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	if data.Inner == nil || (*data.Inner).Field != "test" {
		t.Fatalf("incorrect value. expected=test actual=%v", data.Inner)
	}
}

func TestStringPointer_MultiplePointersFromDefault_Ok(t *testing.T) {
	data := struct {
		Field ***string `default:"test"`
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	if data.Field == nil || ***data.Field != "test" {
		t.Fatalf("incorrect value. expected=test actual=%s", ***data.Field)
	}
}

func TestStructPointer_MultiplePointer_Ok(t *testing.T) {
	data := struct {
		Inner ***struct {
			Field string `default:"test"`
		}
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	if data.Inner == nil || (***data.Inner).Field != "test" {
		t.Fatalf("incorrect value. expected=test actual=%v", data.Inner)
	}
}

func TestStringPoint_WithoutValue_Ok(t *testing.T) {
	data := struct {
		Field *string
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	if data.Field != nil {
		t.Fatal("field not nil")
	}
}

func TestStructPoint_WithoutValue_Ok(t *testing.T) {
	data := struct {
		Inner *struct {
			Field string
		}
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	if data.Inner != nil {
		t.Fatalf("incorrect value. expected=test actual=%v", data.Inner)
	}
}

func TestStructPoint_FewInnerOneWithValue_Ok(t *testing.T) {
	data := struct {
		Inner1 *struct {
			Field string
		}
		Inner2 struct {
			Inner21 **struct {
				Field string
			}
			Field string
		}
		Inner3 *struct {
			Inner31 *struct {
				Field string `default:"test"`
			}
		}
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	if data.Inner1 != nil {
		t.Fatalf("incorrect value. expected=nil actual=%v", data.Inner1)
	}
	if data.Inner2.Inner21 != nil {
		t.Fatalf("incorrect value. expected=nil actual=%v", data.Inner2.Inner21)
	}
	if data.Inner3 == nil || data.Inner3.Inner31 == nil || data.Inner3.Inner31.Field != "test" {
		t.Fatalf("incorrect value. expected=not_nil actual=%#v", data)
	}
}
