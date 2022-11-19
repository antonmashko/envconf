package envconf

import (
	"reflect"
)

type field interface {
	Init() error
	Define() error
}

type emptyField struct{}

func (emptyField) Init() error {
	return nil
}

func (emptyField) Define() error {
	return nil
}

func createFieldFromValue(val reflect.Value, parent *structType, tag reflect.StructField) field {
	switch val.Kind() {
	case reflect.Struct:
		return newStructType(val, parent, tag)
	case reflect.Ptr:
		return newPtrType(val, parent, tag)
	case reflect.Interface:
		// in development
		return &interfaceType{}
	case reflect.Map, reflect.Slice, reflect.Array:
		// in development
		return &collectionType{}
	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Uintptr:
		// unsupported types
		return emptyField{}
	default:
		return newPrimitiveType(val, parent, tag)
	}
}
