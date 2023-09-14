package envconf

import (
	"net/url"
	"reflect"
	"time"
)

const fieldNameDelim = "."

type field interface {
	name() string
	parent() field
	init() error
	define() error
	isSet() bool
	structField() reflect.StructField
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

func createFieldFromValue(v reflect.Value, p *structType, t reflect.StructField) field {
	switch v.Kind() {
	case reflect.Struct:
		switch v.Interface().(type) {
		case url.URL, time.Time:
			return newPrimitiveType(v, p, t)
		default:
			return newStructType(v, p, t)
		}
	case reflect.Ptr:
		return newPtrType(v, p, t)
	case reflect.Interface:
		return emptyField{}
	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Uintptr:
		// unsupported types
		return emptyField{}
	default:
		return newPrimitiveType(v, p, t)
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
