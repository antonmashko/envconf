package envconf

import (
	"errors"
	"reflect"

	"github.com/antonmashko/envconf/external"
)

type structType struct {
	*configField

	sname string
	v     reflect.Value
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

	s := newStructType(v, newConfigField(nil, reflect.StructField{}, parser))
	return s, nil
}

func newStructType(val reflect.Value, f *configField) *structType {
	sname := f.Tag.Get("envconf")
	return &structType{
		sname:       sname,
		configField: f,
		v:           val,
		ext:         external.NilContainer{},
		fields:      make([]field, val.NumField()),
	}
}

func (s *structType) name() string {
	if s.sname != "" {
		return s.sname
	}
	return s.Name
}

func (s *structType) isSet() bool {
	return s.hasValue
}

func (s *structType) externalSource() external.ExternalSource {
	return s.ext
}

func (s *structType) init() error {
	s.fields = make([]field, s.v.NumField())
	rt := s.v.Type()
	for i := 0; i < s.v.NumField(); i++ {
		rfield := s.v.Field(i)
		f := createFieldFromValue(rfield, newConfigField(s, rt.Field(i), s.parser))
		if err := f.init(); err != nil {
			return err
		}
		s.parser.fieldInitialized(f)
		s.fields[i] = f
	}
	return nil
}

func (s *structType) define() error {
	if s.parentField != nil {
		s.ext = external.AsExternalSource(s.Name, s.parentField.externalSource())
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
					FieldName: fullname(f, fieldNameDelim),
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
