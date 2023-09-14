package envconf

import (
	"reflect"
)

type ptrType struct {
	field

	p  *structType
	v  reflect.Value
	sf reflect.StructField

	tmp *reflect.Value
}

func newPtrType(v reflect.Value, p *structType, sf reflect.StructField) *ptrType {
	return &ptrType{
		field: emptyField{},
		p:     p,
		v:     v,
		sf:    sf,
		tmp:   nil,
	}
}

func (f *ptrType) init() error {
	tmp := f.v
	for tmp.Kind() == reflect.Ptr {
		if tmp.IsNil() {
			nv := reflect.New(tmp.Type().Elem())
			if f.tmp == nil {
				f.tmp = &nv
			} else {
				tmp.Set(nv)
			}
			tmp = nv
		} else {
			tmp = tmp.Elem()
		}
	}

	f.field = createFieldFromValue(tmp, f.p, f.sf)
	return f.field.init()
}

func (f *ptrType) define() error {
	err := f.field.define()
	if err != nil {
		return err
	}

	if f.field.isSet() && f.tmp != nil {
		f.v.Set(*f.tmp)
	}

	return nil
}

func (f *ptrType) name() string {
	return f.sf.Name
}

func (f *ptrType) parent() field {
	if f.p == nil {
		return nil
	}
	return f.p
}

func (f *ptrType) structField() reflect.StructField {
	return f.sf
}
