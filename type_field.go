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
	v  reflect.Value
	p  *structType
	sf reflect.StructField

	flag     *flagSource          // flag value
	env      *envSource           // env value
	ext      *externalValueSource // external value
	def      *defaultValueSource  // default value
	required bool                 // if it defined true, value should be defined
	desc     string               // description

	definedValue *definedValue
}

func newFieldType(v reflect.Value, p *structType, sf reflect.StructField) *fieldType {
	desc := sf.Tag.Get(tagDescription)
	required, _ := strconv.ParseBool(sf.Tag.Get(tagRequired))
	f := &fieldType{
		p:        p,
		v:        v,
		sf:       sf,
		def:      newDefaultValueSource(sf),
		required: required,
		desc:     desc,
	}
	f.flag = newFlagSource(f, sf, desc)
	f.env = newEnvSource(f, sf)
	f.ext = newExternalValueSource(f, p.parser.external)
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

func (t *fieldType) define() error {
	// create correct parse priority
	priority := t.p.parser.PriorityOrder()
	for _, p := range priority {
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
		if !ok {
			continue
		}

		var err error
		if str, ok := v.(string); ok {
			err = setFromString(t.v, str)
		} else {
			err = setFromInterface(t.v, v)
		}

		if err != nil {
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

	return ErrConfigurationNotFound
}

func setFromInterface(field reflect.Value, value interface{}) error {
	ival := reflect.ValueOf(value)
	itype := ival.Type()
	if field.Type() == itype {
		field.Set(ival)
		return nil
	}

	switch field.Kind() {
	case reflect.Struct:
		if itype.Kind() != reflect.Map {
			return ErrUnsupportedType
		}
		iter := ival.MapRange()
		for iter.Next() {
			if err := setFromInterface(field.FieldByName(iter.Key().String()), iter.Value()); err != nil {
				return err
			}
		}
		return nil
	case reflect.Array:
		if ikind := itype.Kind(); ikind != reflect.Array && ikind != reflect.Slice {
			return fmt.Errorf("unable to cast %s to array", itype)
		}
		length := ival.Len()
		for i := 0; i < length; i++ {
			if err := setFromInterface(field.Index(i), ival.Index(i).Interface()); err != nil {
				return err
			}
		}
		return nil
	case reflect.Slice:
		if ikind := itype.Kind(); ikind != reflect.Array && ikind != reflect.Slice {
			return fmt.Errorf("unable to cast %s to slice", itype)
		}
		length := ival.Len()
		vtype := field.Type()
		rsl := reflect.MakeSlice(vtype, ival.Cap(), length)
		for i := 0; i < length; i++ {
			if err := setFromInterface(rsl.Index(i), ival.Index(i).Interface()); err != nil {
				return err
			}
		}
		field.Set(rsl)
		return nil
	case reflect.Map:
		if itype.Kind() != reflect.Map {
			return fmt.Errorf("unable to cast %s to array", itype)
		}
		ftype := field.Type()
		rmp := reflect.MakeMap(ftype)
		key := ftype.Key()
		elem := ftype.Elem()
		iter := ival.MapRange()
		for iter.Next() {
			rvkey := reflect.New(key).Elem()
			if err := setFromInterface(rvkey, iter.Key().Interface()); err != nil {
				return err
			}
			rvval := reflect.New(elem).Elem()
			if err := setFromInterface(rvval, iter.Value().Interface()); err != nil {
				return err
			}
			rmp.SetMapIndex(rvkey, rvval)
		}
		field.Set(rmp)
		return nil
	default:
		return setFromString(field, fmt.Sprint(value))
	}
}

func setFromString(field reflect.Value, value string) error {
	if implF := asImpl(field); implF != nil {
		return implF([]byte(value))
	}
	oval := value
	value = strings.Trim(value, " ")
	if !field.CanSet() {
		return errors.New("reflect: cannot set")
	}

	// native complex types
	switch field.Interface().(type) {
	case time.Duration:
		d, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		field.SetInt(d.Nanoseconds())
		return nil
	}

	// primitives and collections
	switch field.Kind() {
	case reflect.String:
		field.SetString(oval)
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
		ftype := field.Type()
		key := ftype.Key()
		elem := ftype.Elem()
		for i := range sl {
			idx := strings.IndexRune(sl[i], ':')
			rvkey := reflect.New(key).Elem()
			rvval := reflect.New(elem).Elem()
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
