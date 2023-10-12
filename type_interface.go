package envconf

import (
	"reflect"

	"github.com/antonmashko/envconf/external"
)

type interfaceType struct {
	*configField
	f field // underline type
	v reflect.Value
}

func newInterfaceType(v reflect.Value, f *configField) *interfaceType {
	return &interfaceType{
		f:           emptyField{},
		configField: f,
		v:           v,
	}
}

func (i *interfaceType) externalSource() external.ExternalSource {
	return i.f.externalSource()
}

func (i *interfaceType) init() error {
	return i.f.init()
}

func (i *interfaceType) define() error {
	if !i.v.IsNil() {
		i.f = createFieldFromValue(i.v.Elem(), i.configField)
	} else {
		i.f = &interfaceFieldType{
			fieldType: newFieldType(i.v, i.configField),
		}
	}
	if err := i.f.init(); err != nil {
		return err
	}
	return i.f.define()
}
