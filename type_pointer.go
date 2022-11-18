package envconf

import (
	"reflect"
)

type ptrType struct {
	field
	v   reflect.Value
	tmp *reflect.Value
	tag reflect.StructField
}

func (f *ptrType) Init(val reflect.Value, parent *structType, tag reflect.StructField) error {
	f.v = val
	f.tag = tag

	tmp := val
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

	f.field = createFieldFromValue(tmp)
	return f.field.Init(tmp, parent, tag)
}

func (f *ptrType) Define() error {
	err := f.field.Define()
	if err != nil {
		return err
	}

	if f.tmp != nil {
		// FIXME: field should not be initialized if none of the fields has a value
		f.v.Set(*f.tmp)
	}

	return nil
}
