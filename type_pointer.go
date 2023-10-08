package envconf

import (
	"reflect"

	"github.com/antonmashko/envconf/external"
)

type ptrType struct {
	*configField

	f        field //pointer value
	v        reflect.Value
	hasValue bool

	tmp *reflect.Value
}

func newPtrType(v reflect.Value, f *configField) *ptrType {
	return &ptrType{
		f:           emptyField{},
		configField: f,
		v:           v,
		tmp:         nil,
	}
}

func (p *ptrType) externalSource() external.ExternalSource {
	return p.f.externalSource()
}

func (p *ptrType) isSet() bool {
	return p.hasValue
}

func (p *ptrType) init() error {
	tmp := p.v
	for tmp.Kind() == reflect.Ptr {
		if tmp.IsNil() {
			nv := reflect.New(tmp.Type().Elem())
			if p.tmp == nil {
				p.tmp = &nv
			} else {
				tmp.Set(nv)
			}
			tmp = nv
		} else {
			tmp = tmp.Elem()
		}
	}

	p.f = createFieldFromValue(tmp, p.configField)
	return p.f.init()
}

func (p *ptrType) define() error {
	err := p.f.define()
	if err != nil {
		return err
	}

	if p.f.isSet() && p.tmp != nil {
		if !p.v.CanSet() {
			return nil
		}
		p.v.Set(*p.tmp)
		p.hasValue = true
	}

	return nil
}
