package envconf

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type definedValue struct {
	source ConfigSource
	value  interface{}
}

type primitiveType struct {
	v  reflect.Value
	p  *structType
	sf reflect.StructField

	flag     Var    // flag value
	env      Var    // env value
	def      Var    // default value
	required bool   // if it defined true, value should be defined
	desc     string // description

	definedValue *definedValue
}

func newPrimitiveType(v reflect.Value, p *structType, sf reflect.StructField) *primitiveType {
	desc := sf.Tag.Get(tagDescription)
	required, _ := strconv.ParseBool(sf.Tag.Get(tagRequired))
	return &primitiveType{
		p:        p,
		v:        v,
		sf:       sf,
		flag:     newFlagSource(sf, desc),
		env:      newEnvSource(sf),
		def:      newDefaultValueSource(sf),
		required: required,
		desc:     desc,
	}
}

func (t *primitiveType) name() string {
	return t.sf.Name
}

func (t *primitiveType) parent() field {
	return t.p
}

func (t *primitiveType) init() error {
	return nil
}

func (t *primitiveType) define() error {
	// validate reflect value
	if !t.v.IsValid() {
		return errInvalidFiled
	}
	if !t.v.CanSet() {
		return fmt.Errorf("%s: %w", t.Name(), errFiledIsNotSettable)
	}

	// create correct parse priority
	priority := t.p.parser.PriorityOrder()
	for _, p := range priority {
		var v Var
		switch p {
		case FlagVariable:
			v = t.flag
		case EnvVariable:
			v = t.env
		case ExternalSource:
			values := []Value{t}
			var parent Value = t.p
			for parent != nil && parent.Name() != "" {
				values = append([]Value{parent}, values...)
				parent = parent.Owner()
			}
			value, exists := t.p.parser.external.Get(values...)
			if exists {
				t.v.Set(reflect.ValueOf(value))
				return nil
			}
			continue
		case DefaultValue:
			v = t.def
		}

		if str, ok := v.Value(); ok {
			t.definedValue = &definedValue{
				source: p,
				value:  v,
			}
			// set value
			return setFromString(t.v, str)
		}
	}

	return errConfigurationNotSpecified
}

func (t *primitiveType) Owner() Value {
	return t.p
}

func (t *primitiveType) Name() string {
	return t.name()
}

func (t *primitiveType) Tag() reflect.StructField {
	return t.sf
}

func (t *primitiveType) isSet() bool {
	return t.definedValue != nil
}

func (t *primitiveType) IsRequired() bool {
	return t.required
}

func setFromString(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Bool:
		i, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		var err error
		if _, ok := field.Interface().(time.Duration); ok {
			var d time.Duration
			d, err = time.ParseDuration(value)
			if err != nil {
				return err
			}
			i = d.Nanoseconds()
		} else {
			i, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(i)
	case reflect.Complex64:
		i, err := strconv.ParseComplex(value, 64)
		if err != nil {
			return err
		}
		field.SetComplex(i)
	case reflect.Complex128:
		i, err := strconv.ParseComplex(value, 128)
		if err != nil {
			return err
		}
		field.SetComplex(i)
	case reflect.String:
		field.SetString(value)
	default:
		return ErrUnsupportedType
	}
	return nil
}
