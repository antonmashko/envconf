package envconf

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/antonmashko/envconf/option"
)

type valueExtractor interface {
	Value() (interface{}, bool)
}

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

func newFieldType(v reflect.Value, p field, sf reflect.StructField, ext *externalConfig, parseOrder []option.ConfigSource) *fieldType {
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
	f.ext = newExternalValueSource(f, ext)
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

func (t *fieldType) init() error {
	return nil
}

func (t *fieldType) readConfigValue() (interface{}, option.ConfigSource, error) {
	// create correct parse priority
	for _, p := range t.parseOrder {
		var vr valueExtractor
		switch p {
		case option.FlagVariable:
			vr = t.flag
		case option.EnvVariable:
			vr = t.env
		case option.ExternalSource:
			vr = t.ext
		case option.DefaultValue:
			vr = t.def
		}

		v, ok := vr.Value()
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
	v, p, err := t.readConfigValue()
	if str, ok := v.(string); ok && p != option.ExternalSource {
		v, err = setFromString(t.v, str)
	}

	if err != nil {
		if err == ErrConfigurationNotFound {
			return err
		}
		return &Error{
			Inner:     fmt.Errorf("type=%s source=%s. %w", t.sf.Type, p, err),
			FieldName: fullname(t),
			Message:   "cannot set",
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
		res := []byte(value)
		return res, implF(res)
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
