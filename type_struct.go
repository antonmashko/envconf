package envconf

import (
	"reflect"
)

var ErrUnsupportedType = errUnsupportedType

type structType struct {
	parser *EnvConf
	parent *structType

	v   reflect.Value
	t   reflect.Type
	tag reflect.StructField

	fields []field
}

func (s *structType) Init(val reflect.Value, parent *structType, tag reflect.StructField) error {
	s.parent = parent
	s.v = val
	s.t = val.Type()
	s.tag = tag
	s.fields = make([]field, s.v.NumField())

	for i := 0; i < s.v.NumField(); i++ {
		rfield := s.v.Field(i)
		f := createFieldFromValue(rfield)
		if err := f.Init(rfield, s, s.t.Field(i)); err != nil {
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
