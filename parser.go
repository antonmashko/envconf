package envconf

import (
	"errors"
	"flag"
	"reflect"
)

const Separator = "."

var (
	ErrNilData = errors.New("nil data")
)

//Parse fiend tag(annotations) for each field as set value
func Parse(data interface{}) error {
	if data == nil {
		return ErrNilData
	}
	p, err := newParser(data)
	if err != nil {
		return err
	}
	if err = p.Init(); err != nil {
		return err
	}
	if UseCustomHelp {
		flag.Usage = (&help{p: p}).usage
	}
	flag.Parse()
	return p.Parse()
}

type parser struct {
	value  reflect.Value
	rtype  reflect.Type
	parent *parser
	childs []*parser
	values []*value
}

func newParser(data interface{}) (*parser, error) {
	return newChildParser(nil, reflect.ValueOf(data))
}

func newChildParser(parent *parser, rvalue reflect.Value) (*parser, error) {
	p := &parser{}
	p.value = rvalue
	if p.value.Kind() == reflect.Ptr {
		// check on nil
		if p.value.IsNil() {
			return nil, errors.New("Data is nil pointer")
		}

		p.value = p.value.Elem() // get value from pointer
	}
	p.rtype = p.value.Type() // remember type
	p.childs = make([]*parser, 0)
	p.values = make([]*value, 0)
	p.parent = parent
	return p, nil
}

func (p *parser) Init() error {
	for i := 0; i < p.value.NumField(); i++ {
		v := p.value.Field(i)
		if v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr {
			cp, err := newChildParser(p, v)
			if err != nil {
				return err
			}
			p.childs = append(p.childs, cp)
			err = cp.Init()
			if err != nil {
				return err
			}
			continue
		}
		//TODO: check on another type
		vl := newValue(v, p.rtype.Field(i))
		vl.owner = p
		p.values = append(p.values, vl)
	}
	return nil
}

func (p *parser) Parse() error {
	for _, v := range p.values {
		err := v.define()
		if err != nil {
			return err
		}
	}
	for _, v := range p.childs {
		err := v.Parse()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) Path() string {
	if p.parent == nil {
		return ""
	}
	return p.parent.Path() + Separator + p.rtype.Name()
}
