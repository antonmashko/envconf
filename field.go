package envconf

import (
	"log"
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

func createFieldFromValue(v reflect.Value, p *structType, t reflect.StructField) field {
	log.Printf("creating type--> %s. Name:%s %p %p %v", v.Kind(), t.Name, p, p.parent, p.parser)
	switch v.Kind() {
	case reflect.Struct:
		return newStructType(v, p, t)
	case reflect.Ptr:
		return newPtrType(v, p, t)
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
		return newPrimitiveType(v, p, t)
	}
}
