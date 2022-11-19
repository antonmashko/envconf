package envconf

import (
	"errors"
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

func newParentStructType(data interface{}, parser *EnvConf) (*structType, error) {
	v := reflect.ValueOf(data)
	for v.Kind() == reflect.Ptr {
		// check on nil
		if v.IsNil() {
			return nil, ErrNilData
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("invalid type")
	}

	s := newStructType(v, nil, reflect.StructField{})
	s.parser = parser
	return s, nil
}

func newStructType(val reflect.Value, parent *structType, tag reflect.StructField) *structType {
	var p *EnvConf
	if parent != nil {
		p = parent.parser
	}
	return &structType{
		parent: parent,
		parser: p,
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
		s.parser.fieldInitialized(f)
		s.fields[i] = f
	}
	return nil
}

func (s *structType) Define() error {
	for _, f := range s.fields {
		err := f.Define()
		if err != nil {
			if rf, ok := f.(requiredField); ok && rf.IsRequired() {
				return &Error{
					Message:   "failed to define field",
					Inner:     err,
					FieldName: fullname(f.(Value)),
				}
			}

			s.parser.fieldNotDefined(f, err)
			if err == errConfigurationNotSpecified {
				continue
			}
			return err
		}
		s.parser.fieldDefined(f)
	}
	return nil
}

func (s *structType) Owner() Value {
	if s.parent == nil {
		return nil
	}
	return s.parent
}

func (s *structType) Name() string {
	return s.tag.Name
}

func (s *structType) Tag() reflect.StructField {
	return s.tag
}
