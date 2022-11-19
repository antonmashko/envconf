package envconf

import (
	"reflect"
)

var ErrUnsupportedType = errUnsupportedType

type structType struct {
	parent *structType

	v   reflect.Value
	t   reflect.Type
	tag reflect.StructField

	fields []field
}

func newStructType(val reflect.Value, parent *structType, tag reflect.StructField) *structType {
	return &structType{
		parent: parent,
		v:      val,
		t:      val.Type(),
		tag:    tag,
		fields: make([]field, val.NumField()),
	}
}

func (s *structType) Init() error {
	s.fields = make([]field, s.v.NumField())
	for i := 0; i < s.v.NumField(); i++ {
		rfield := s.v.Field(i)
		f := createFieldFromValue(rfield, s, s.t.Field(i))
		if err := f.Init(); err != nil {
			return err
		}
		s.fields[i] = f
	}
	return nil
}

func (s *structType) Define() error {
	for _, f := range s.fields {
		if err := f.Define(); err != nil {
			return err
		}
	}
	return nil
}
