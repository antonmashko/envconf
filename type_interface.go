package envconf

import (
	"reflect"
)

type interfaceType struct {
	field  // underline type
	v      reflect.Value
	p      field
	sf     reflect.StructField
	parser *EnvConf
}

func newInterfaceType(v reflect.Value, p field, sf reflect.StructField, parser *EnvConf) *interfaceType {
	return &interfaceType{
		field:  emptyField{},
		v:      v,
		p:      p,
		sf:     sf,
		parser: parser,
	}
}

func (t *interfaceType) name() string {
	return t.sf.Name
}

func (t *interfaceType) parent() field {
	return t.p
}

func (t *interfaceType) define() error {
	if t.v.IsValid() && !t.v.IsZero() {
		t.field = createFieldFromValue(t.v.Elem(), t.p, t.sf, t.parser)
		if err := t.field.init(); err != nil {
			return err
		}
	}
	return t.field.define()
}

func (t *interfaceType) structField() reflect.StructField {
	return t.sf
}
