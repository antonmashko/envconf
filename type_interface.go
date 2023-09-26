package envconf

import "reflect"

type interfaceType struct {
	fv     field // underline type
	v      reflect.Value
	p      field
	sf     reflect.StructField
	parser *EnvConf
}

func newInterfaceType(v reflect.Value, p field, sf reflect.StructField, parser *EnvConf) *interfaceType {
	return &interfaceType{
		fv:     emptyField{},
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

func (t *interfaceType) init() error {
	return t.fv.init()
}

func (t *interfaceType) define() error {
	if t.v.IsValid() && !t.v.IsZero() {
		t.fv = createFieldFromValue(t.v.Elem(), t.p, t.sf, t.parser)
		if err := t.fv.init(); err != nil {
			return err
		}
	}
	return t.fv.define()
}

func (t *interfaceType) isSet() bool {
	return t.fv.isSet()
}

func (t *interfaceType) structField() reflect.StructField {
	return t.sf
}
