package envconf

import "reflect"

type interfaceType struct{}

func (t *interfaceType) Init(val reflect.Value, parent *structType, tag reflect.StructField) error {
	// not supported
	return nil
}

func (t *interfaceType) Define() error {
	return nil
}
