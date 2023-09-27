package envconf

import (
	"errors"
	"reflect"

	"github.com/antonmashko/envconf/external"
)

type structType struct {
	parser *EnvConf
	p      field

	sname string
	v     reflect.Value
	t     reflect.Type
	tag   reflect.StructField
	ext   external.ExternalSource

	hasValue bool
	fields   []field
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

	s := newStructType(v, nil, reflect.StructField{}, parser)
	return s, nil
}

func newStructType(val reflect.Value, parent field, tag reflect.StructField, parser *EnvConf) *structType {
	sname := tag.Tag.Get("envconf")
	return &structType{
		parser: parser,
		sname:  sname,
		p:      parent,
		v:      val,
		t:      val.Type(),
		tag:    tag,
		ext:    external.NilContainer{},
		fields: make([]field, val.NumField()),
	}
}

func (s *structType) name() string {
	if s.sname != "" {
		return s.sname
	}
	return s.tag.Name
}

func (s *structType) parent() field {
	if s.p == nil {
		return nil
	}
	return s.p
}

func (s *structType) structField() reflect.StructField {
	return s.tag
}

func (s *structType) isSet() bool {
	return s.hasValue
}

func (s *structType) externalSource() external.ExternalSource {
	return s.ext
}

func (s *structType) init() error {
	s.fields = make([]field, s.v.NumField())
	for i := 0; i < s.v.NumField(); i++ {
		rfield := s.v.Field(i)
		f := createFieldFromValue(rfield, s, s.t.Field(i), s.parser)
		if err := f.init(); err != nil {
			return err
		}
		s.parser.fieldInitialized(f)
		s.fields[i] = f
	}
	return nil
}

func (s *structType) define() error {
	if s.p != nil {
		s.ext = external.AsExternalSource(s.tag.Name, s.p.externalSource())
	}
	for _, f := range s.fields {
		err := f.define()
		if err != nil {
			s.parser.fieldNotDefined(f, err)
			if !errors.Is(err, ErrConfigurationNotFound) {
				return err
			}
			if rf, ok := f.(requiredField); ok && rf.IsRequired() {
				return &Error{
					Message:   "failed to define field",
					Inner:     err,
					FieldName: fullname(f),
				}
			}
		}

		if f.isSet() {
			s.hasValue = true
			s.parser.fieldDefined(f)
		}
	}
	return nil
}
