package envconf

import (
	"encoding"
	"reflect"

	"github.com/antonmashko/envconf/external"
)

const fieldNameDelim = "."

type namedField interface {
	name() string
	parent() field
}

type field interface {
	namedField
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

func createFieldFromValue(v reflect.Value, f *configField) field {
	if v.Kind() == reflect.Pointer {
		return newPtrType(v, f)
	}
	// validate reflect value
	if !v.CanInterface() {
		return emptyField{}
	}

	// implementations check
	implF := asImpl(v)
	if implF != nil {
		return &customSetFieldType{newFieldType(v, f)}
	}

	switch v.Kind() {
	case reflect.Struct:
		return newStructType(v, f)
	case reflect.Interface:
		if v.IsNil() {
			return &interfaceFieldType{
				fieldType: newFieldType(v, f),
			}
		}
		return newInterfaceType(v, f)
	case reflect.Array, reflect.Slice:
		return &collectionSliceType{
			collectionType: newCollectionType(v, f),
		}
	case reflect.Map:
		return &collectionMapType{
			collectionType: newCollectionType(v, f),
		}
	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Uintptr:
		// unsupported types
		return emptyField{}
	default:
		return newFieldType(v, f)
	}
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

func asConfigField(f field) *configField {
	switch ft := f.(type) {
	case *fieldType:
		return ft.configField
	case *ptrType:
		return asConfigField(ft.f)
	case *interfaceType:
		return asConfigField(ft.f)
	case *interfaceFieldType:
		return ft.configField
	case *customSetFieldType:
		return ft.configField
	case *collectionType:
		return ft.configField
	default:
		return nil
	}
}

func fullname(f namedField, delim string) string {
	if f == nil {
		return ""
	}
	name := f.name()
	for {
		f = f.parent()
		if f == nil {
			break
		}
		pname := f.name()
		if pname == "" {
			break
		}
		name = pname + delim + name
	}
	return name
}
