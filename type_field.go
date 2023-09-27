package envconf

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/antonmashko/envconf/external"
	"github.com/antonmashko/envconf/option"
)

type definedValue struct {
	source option.ConfigSource
	value  interface{}
}

type fieldType struct {
	v          reflect.Value
	p          field
	sf         reflect.StructField
	parseOrder []option.ConfigSource

	flag     *flagSource          // flag value
	env      *envSource           // env value
	ext      *externalValueSource // external value
	def      *defaultValueSource  // default value
	required bool                 // if it defined true, value should be defined
	desc     string               // description

	definedValue *definedValue
}

func newFieldType(v reflect.Value, p field, sf reflect.StructField, parseOrder []option.ConfigSource, allowEnvInjection bool) *fieldType {
	desc := sf.Tag.Get(tagDescription)
	required, _ := strconv.ParseBool(sf.Tag.Get(tagRequired))
	f := &fieldType{
		v:          v,
		p:          p,
		sf:         sf,
		parseOrder: parseOrder,
		def:        newDefaultValueSource(sf),
		required:   required,
		desc:       desc,
	}
	f.flag = newFlagSource(f, sf, desc)
	f.env = newEnvSource(f, sf)
	f.ext = newExternalValueSource(f, allowEnvInjection)
	return f
}

func (t *fieldType) name() string {
	return t.sf.Name
}

func (t *fieldType) parent() field {
	return t.p
}

func (t *fieldType) isSet() bool {
	return t.definedValue != nil
}

func (t *fieldType) structField() reflect.StructField {
	return t.sf
}

func (t *fieldType) IsRequired() bool {
	return t.required
}

func (t *fieldType) externalSource() external.ExternalSource {
	return external.NilContainer{}
}

func (t *fieldType) init() error {
	return nil
}

func (t *fieldType) readValue() (interface{}, option.ConfigSource, error) {
	// create correct parse priority
	for _, p := range t.parseOrder {
		var v interface{}
		var ok bool
		switch p {
		case option.FlagVariable:
			v, ok = t.flag.Value()
		case option.EnvVariable:
			v, ok = t.env.Value()
		case option.ExternalSource:
			v, p, ok = t.ext.Value()
		case option.DefaultValue:
			v, ok = t.def.Value()
		}
		if ok {
			return v, p, nil
		}
	}
	return nil, option.ConfigSource(-1), ErrConfigurationNotFound
}

func (t *fieldType) define() error {
	if t.definedValue != nil {
		return nil
	}
	v, p, err := t.readValue()
	if err != nil {
		return err
	}

	if str, ok := v.(string); ok && p != option.ExternalSource {
		v, err = setFromString(t.v, str)
		if err != nil {
			return &Error{
				Inner:     fmt.Errorf("type=%s source=%s. %w", t.sf.Type, p, err),
				FieldName: fullname(t),
				Message:   "cannot set",
			}
		}
	}

	t.definedValue = &definedValue{
		source: p,
		value:  v,
	}
	return nil
}

func (t *fieldType) defineFromValue(v interface{}, p option.ConfigSource) error {
	str, ok := v.(string)
	if ok {
		var err error
		v, err = setFromString(t.v, str)
		if err != nil {
			return &Error{
				Inner:     fmt.Errorf("type=%s source=%s. %w", t.sf.Type, p, err),
				FieldName: fullname(t),
				Message:   "cannot set",
			}
		}
	}

	t.definedValue = &definedValue{
		source: p,
		value:  v,
	}

	return nil
}

func setFromString(field reflect.Value, value string) (interface{}, error) {
	if implF := asImpl(field); implF != nil {
		return value, implF([]byte(value))
	}
	oval := value
	value = strings.Trim(value, " ")
	if !field.CanSet() {
		return nil, errors.New("reflect: cannot set")
	}

	// native complex types
	switch field.Interface().(type) {
	case time.Duration:
		d, err := time.ParseDuration(value)
		if err != nil {
			return nil, err
		}
		field.SetInt(d.Nanoseconds())
		return d, nil
	}

	// primitives and collections
	switch field.Kind() {
	case reflect.String:
		field.SetString(oval)
		return oval, nil
	case reflect.Bool:
		i, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		field.SetBool(i)
		return i, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 0, field.Type().Bits())
		if err != nil {
			return nil, err
		}
		field.SetInt(i)
		return i, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 0, field.Type().Bits())
		if err != nil {
			return nil, err
		}
		field.SetUint(i)
		return i, nil
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(value, field.Type().Bits())
		if err != nil {
			return nil, err
		}
		field.SetFloat(i)
		return i, nil
	case reflect.Complex64, reflect.Complex128:
		i, err := strconv.ParseComplex(value, field.Type().Bits())
		if err != nil {
			return nil, err
		}
		field.SetComplex(i)
		return i, nil
	default:
		return nil, ErrUnsupportedType
	}
}
