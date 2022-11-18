package envconf

import (
	"fmt"
	"log"
	"reflect"
)

var ErrUnsupportedType = errUnsupportedType

type structFieldData struct {
	parent *structFieldData

	v   reflect.Value
	t   reflect.Type
	tag reflect.StructField

	values map[string]Value
}

func newStructFieldData(val reflect.Value, parent *structFieldData, tag reflect.StructField) *structFieldData {
	t := val.Type()
	log.Println(tag.Tag)
	return &structFieldData{
		parent: parent,
		v:      val,
		t:      t,
		tag:    tag,
		values: make(map[string]Value),
	}
}

func (s *structFieldData) Serialize(data string) error {
	for i := 0; i < s.v.NumField(); i++ {
		field := s.v.Field(i)
		kind := field.Kind()

		// using loop for handling case with more than one referencing
		for kind == reflect.Ptr {
			if field.IsNil() {
				// FIXME: field should not be initialized if none of the fields has a value
				v := reflect.New(field.Type().Elem())
				field.Set(v)
				continue
			} else {
				field = field.Elem()
			}
			kind = field.Kind()
		}

		tag := s.t.Field(i)

		switch kind {
		case reflect.Struct:
			innerStruct := newStructFieldData(field, s, s.t.Field(i))
			return innerStruct.Serialize("")
		case reflect.Interface:
			continue
		default:
			v := newValue(&parser{
				external: &emptyExt{},
			}, field, tag)
			if err := v.define(); err != nil && err != errConfigurationNotSpecified {
				return fmt.Errorf("%s: %w", tag.Name, err)
			}
			// define return `errConfigurationNotSpecified` error if no configuration for that field
			s.values[field.String()] = v
		}
	}
	return nil
}
