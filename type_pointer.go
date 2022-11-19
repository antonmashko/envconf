package envconf

import (
	"reflect"
)

type ptrType struct {
	field
	parent *structType
	v      reflect.Value
	tmp    *reflect.Value
	tag    reflect.StructField
}

func newPtrType(val reflect.Value, parent *structType, tag reflect.StructField) *ptrType {
	return &ptrType{
		field:  emptyField{},
		parent: parent,
		v:      val,
		tag:    tag,
		tmp:    nil,
	}
}

func (f *ptrType) Init() error {
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

	f.field = createFieldFromValue(tmp, f.parent, f.tag)
	return f.field.Init()
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
