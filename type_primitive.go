package envconf

import (
	"encoding"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
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
	f := &primitiveType{
		p:        p,
		v:        v,
		sf:       sf,
		def:      newDefaultValueSource(sf),
		required: required,
		desc:     desc,
	}
	f.flag = newFlagSource(f, sf, desc)
	f.env = newEnvSource(f, sf)
	return f
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
			_, ok := t.p.parser.external.Get(values...)
			if ok {
				// field defined in external source
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
	value = strings.Trim(value, " ")
	// native complex types
	switch it := field.Interface().(type) {
	case encoding.TextUnmarshaler:
		return it.UnmarshalText([]byte(value))
	case encoding.BinaryUnmarshaler:
		return it.UnmarshalBinary([]byte(value))
	case url.URL:
		url, err := url.Parse(value)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(*url))
		return nil
	case time.Duration:
		d, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		field.SetInt(d.Nanoseconds())
		return nil
	case net.IP:
		field.Set(reflect.ValueOf(net.ParseIP(value)))
		return nil
	case time.Time:
		dt, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(dt))
		return nil
	}

	// primitives and collections
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		i, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 0, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 0, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(value, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(i)
	case reflect.Complex64, reflect.Complex128:
		i, err := strconv.ParseComplex(value, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetComplex(i)
	case reflect.Array:
		sl := strings.Split(value, ",")
		for i := range sl {
			err := setFromString(field.Index(i), sl[i])
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		sl := strings.Split(value, ",")
		rsl := reflect.MakeSlice(field.Type(), len(sl), cap(sl))
		for i := range sl {
			err := setFromString(rsl.Index(i), sl[i])
			if err != nil {
				return err
			}
		}
		field.Set(rsl)
	case reflect.Map:
		sl := strings.Split(value, ",")
		rmp := reflect.MakeMap(field.Type())
		for i := range sl {
			idx := strings.IndexRune(sl[i], ':')
			rvkey := reflect.New(field.Type().Key()).Elem()
			rvval := reflect.New(field.Type().Elem()).Elem()
			if idx == -1 {
				if err := setFromString(rvkey, sl[i]); err != nil {
					return err
				}
			} else {
				if err := setFromString(rvkey, sl[i][:idx]); err != nil {
					return err
				}
				if err := setFromString(rvval, sl[i][idx+1:]); err != nil {
					return err
				}
			}
			rmp.SetMapIndex(rvkey, rvval)
		}
		field.Set(rmp)
	default:
		return ErrUnsupportedType
	}
	return nil
}
