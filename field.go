package envconf

import (
	"encoding"
	"reflect"

	"github.com/antonmashko/envconf/external"
)

const fieldNameDelim = "."

type field interface {
	name() string
	parent() field
	init() error
	define() error
	isSet() bool
	structField() reflect.StructField
	externalSource() external.ExternalSource
}

type requiredField interface {
	IsRequired() bool
}

type emptyField struct{}

func (emptyField) init() error {
	return nil
}

func (emptyField) define() error {
	return nil
}

func (emptyField) isSet() bool {
	return false
}

func (emptyField) parent() field {
	return nil
}

func (emptyField) name() string {
	return ""
}

func (emptyField) structField() reflect.StructField {
	return reflect.StructField{}
}

func (emptyField) externalSource() external.ExternalSource {
	return external.NilContainer{}
}

func createFieldFromValue(v reflect.Value, p field, t reflect.StructField, parser *EnvConf) field {
	if v.Kind() == reflect.Pointer {
		return newPtrType(v, p, t, parser)
	}
	// validate reflect value
	if !v.CanInterface() {
		return emptyField{}
	}

	// implementations check
	implF := asImpl(v)
	if implF != nil {
		return newFieldType(v, p, t, parser.PriorityOrder(), parser.opts.AllowExternalEnvInjection)
	}

	switch v.Kind() {
	case reflect.Struct:
		return newStructType(v, p, t, parser)
	case reflect.Interface:
		return newInterfaceType(v, p, t, parser)
	case reflect.Array, reflect.Slice:
		return &collectionSliceType{
			collectionType: newCollectionType(v, p, t, parser),
		}
	case reflect.Map:
		return &collectionMapType{
			collectionType: newCollectionType(v, p, t, parser),
		}
	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Uintptr:
		// unsupported types
		return emptyField{}
	default:
		return newFieldType(v, p, t, parser.PriorityOrder(), parser.opts.AllowExternalEnvInjection)
	}
}

func fullname(f field) string {
	name := f.name()
	for {
		f = f.parent()
		if f == nil {
			break
		}
		oname := f.name()
		if oname != "" {
			name = oname + fieldNameDelim + name
		}
	}
	return name
}

func asImpl(field reflect.Value) func([]byte) error {
	f := func(v interface{}) func([]byte) error {
		// encoding.TextUnmarshaler
		tu, ok := v.(encoding.TextUnmarshaler)
		if ok {
			return tu.UnmarshalText
		}
		// encoding.BinaryUnmarshaller
		bu, ok := v.(encoding.BinaryUnmarshaler)
		if ok {
			return bu.UnmarshalBinary
		}
		// ----
		return nil
	}
	// NOTE: max double pointer support
	for i := 0; i < 3; i++ {
		resF := f(field.Interface())
		if resF != nil {
			return resF
		}
		if !field.CanAddr() {
			return nil
		}
		field = field.Addr()
		if !field.CanInterface() {
			return nil
		}
	}
	return nil
}

func asFieldType(f field) *fieldType {
	switch ft := f.(type) {
	case *fieldType:
		return ft
	case *ptrType:
		return asFieldType(ft.field)
	case *interfaceType:
		return asFieldType(ft.field)
	default:
		return nil
	}
}
