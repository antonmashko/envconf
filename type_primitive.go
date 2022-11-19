package envconf

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type definedValue struct {
	source Priority
	value  interface{}
}

type primitiveType struct {
	parent   *structType
	v        reflect.Value
	tag      reflect.StructField
	flag     Var    // flag value
	env      Var    // env value
	def      Var    // default value
	required bool   // if it defined true, value should be defined
	desc     string // description

	definedValue *definedValue
}

func newPrimitiveType(val reflect.Value, parent *structType, tag reflect.StructField) *primitiveType {
	desc := tag.Tag.Get(tagDescription)
	required, _ := strconv.ParseBool(tag.Tag.Get(tagRequired))
	return &primitiveType{
		parent:   parent,
		v:        val,
		tag:      tag,
		flag:     newFlagSource(tag, desc),
		env:      newEnvSource(tag),
		def:      newDefaultValueSource(tag),
		required: required,
		desc:     desc,
	}
}

func (t *primitiveType) Init() error {
	return nil
}

func (t *primitiveType) Define() error {
	// validate reflect value
	if !t.v.IsValid() {
		return errInvalidFiled
	}
	if !t.v.CanSet() {
		return fmt.Errorf("%s: %w", t.Name(), errFiledIsNotSettable)
	}

	// create correct parse priority
	priority := priorityOrder()
	for _, p := range priority {
		var v Var
		switch p {
		case FlagPriority:
			v = t.flag
		case EnvPriority:
			v = t.env
		case ExternalPriority:
			values := []Value{t}
			var parent Value = t.parent
			for parent != nil && parent.Name() != "" {
				values = append([]Value{parent}, values...)
				parent = parent.Owner()
			}
			value, exists := t.parent.parser.external.Get(values...)
			if exists {
				t.v.Set(reflect.ValueOf(value))
				return nil
			}
			continue
		case DefaultPriority:
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
	return t.parent
}

func (t *primitiveType) Name() string {
	return t.tag.Name
}

func (t *primitiveType) Tag() reflect.StructField {
	return t.tag
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
	case reflect.String:
		field.SetString(value)
	default:
		return errUnsupportedType
	}
	return nil
}
