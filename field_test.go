package envconf

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEmptyField_Ok(t *testing.T) {
	et := emptyField{}
	if err := et.init(); err != nil {
		t.Fatal("emptyField.init: ", err)
	}
	if err := et.define(); err != nil {
		t.Fatal("emptyField.define: ", err)
	}
	if et.isSet() {
		t.Fatal("emptyField.isSet: true")
	}
	if et.name() != "" {
		t.Fatal("emptyField.name: ", et.name())
	}
	if et.parent() != nil {
		t.Fatal("emptyField.parent: ", et.parent())
	}
	if et.structField().Tag != "" {
		t.Fatal("emptyField.structField: ", et.structField().Tag)
	}
}

func TestCreateField_Ok(t *testing.T) {
	st := struct {
		Field fmt.Stringer
	}{}
	rv := reflect.ValueOf(st)
	rv = rv.Field(0)
	f := createFieldFromValue(rv, newConfigField(emptyField{}, reflect.StructField{Type: rv.Type()}, New()))
	t.Logf("%[1]v %[1]T", f)
}
