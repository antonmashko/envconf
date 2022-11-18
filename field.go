package envconf

import (
	"reflect"
)

type field interface {
	Init(val reflect.Value, parent *structType, tag reflect.StructField) error
	Define() error
}

func createFieldFromValue(v reflect.Value) field {
	switch v.Kind() {
	case reflect.Struct:
		return &structType{}
	case reflect.Ptr:
		return &ptrType{}
	case reflect.Interface:
		return &interfaceType{}
	default:
		return &value{}
	}
}
