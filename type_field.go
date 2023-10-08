package envconf

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/antonmashko/envconf/external"
	"github.com/antonmashko/envconf/option"
)

type fieldType struct {
	*configField
	v  reflect.Value
	sf reflect.StructField
}

func newFieldType(v reflect.Value, s *configField) *fieldType {
	return &fieldType{
		configField: s,
		v:           v,
	}
}

func (*fieldType) externalSource() external.ExternalSource {
	return external.NilContainer{}
}

func (f *fieldType) init() error {
	return f.configField.init(f)
}

func (f *fieldType) define() error {
	v, cs := f.configField.Value()
	if cs == option.NoConfigValue {
		return ErrConfigurationNotFound
	}

	if cs == option.ExternalSource {
		// field should be defined through External.Unmarshal func
		return f.set(v, cs)
	}

	str, ok := v.(string)
	if !ok {
		return &Error{
			Inner:     ErrUnsupportedType,
			FieldName: f.fullName(),
			Message:   "v is not string",
		}
	}
	if !f.v.CanSet() {
		return &Error{
			Inner:     errors.New("reflect: cannot set"),
			FieldName: f.fullName(),
		}
	}

	var err error
	v, err = setFromString(f.v, str)
	if err != nil {
		return &Error{
			Inner:     fmt.Errorf("type=%s source=%s. %w", f.sf.Type, cs, err),
			FieldName: f.fullName(),
			Message:   "cannot set",
		}
	}

	return f.set(v, cs)
}

type interfaceFieldType struct {
	*fieldType
}

func (f *interfaceFieldType) define() error {
	v, cs := f.Value()
	if cs == option.NoConfigValue {
		return ErrConfigurationNotFound
	}
	if cs == option.ExternalSource {
		// field should be defined through External.Unmarshal func
		return f.set(v, cs)
	}
	rv := reflect.ValueOf(v)
	if !rv.Type().AssignableTo(f.v.Type()) {
		return &Error{
			FieldName: f.fullName(),
			Message:   fmt.Sprintf("unable to assign type %s to %s", rv.Type(), f.sf.Type),
		}
	}
	f.v.Set(rv)
	return f.set(v, cs)
}

type customSetFieldType struct {
	*fieldType
}

func (f *customSetFieldType) define() error {
	v, cs := f.Value()
	if cs == option.NoConfigValue {
		return ErrConfigurationNotFound
	}

	if cs == option.ExternalSource {
		// field should be defined through External.Unmarshal func
		return f.set(v, cs)
	}

	str, ok := v.(string)
	if !ok {
		return &Error{
			Inner:     ErrUnsupportedType,
			FieldName: f.fullName(),
			Message:   "v is not string",
		}
	}
	implF := asImpl(f.v)
	if implF == nil {
		return errors.New("setter not found")
	}
	if err := implF([]byte(str)); err != nil {
		return &Error{Inner: err, FieldName: f.fullName()}
	}
	return f.set(v, cs)
}
