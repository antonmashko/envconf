package envconf

import (
	"errors"
	"flag"
	"reflect"

	"github.com/antonmashko/envconf/logger"
)

const Separator = "."

var (
	ErrNilData = errors.New("nil data")
)

var elog = logger.New("envconf", logger.Debug)

func SetLogger(l logger.Logger) {
	if l == nil {
		return
	}
	elog = l
}

//Parse fiend tag(annotations) for each field as set value
func Parse(data interface{}) error {
	return ParseWithExternal(data, nil)
}

func ParseWithExternal(data interface{}, external Config) error {
	if data == nil {
		return ErrNilData
	}
	p, err := newParser(data, external)
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
	if external != nil {
		if err = external.Unmarshal(data); err != nil {
			return err
		}
	}
	return p.Parse()
}

type parser struct {
	value    reflect.Value
	rtype    reflect.Type
	parent   *parser
	external Config
	children []*parser
	values   []*value
}

func newParser(data interface{}, external Config) (*parser, error) {
	return newChildParser(nil, reflect.ValueOf(data), external)
}

func newChildParser(parent *parser, rvalue reflect.Value, external Config) (*parser, error) {
	p := &parser{external: external}
	p.value = rvalue
	if p.value.Kind() == reflect.Ptr {
		// check on nil
		if p.value.IsNil() {
			return nil, ErrNilData
		}
		p.value = p.value.Elem() // get value from pointer
	}
	p.rtype = p.value.Type() // remember type
	p.children = make([]*parser, 0)
	p.values = make([]*value, 0)
	p.parent = parent
	return p, nil
}

func (p *parser) Init() error {
	for i := 0; i < p.value.NumField(); i++ {
		v := p.value.Field(i)
		if v.Kind() == reflect.Struct || v.Kind() == reflect.Ptr {
			cp, err := newChildParser(p, v, p.external)
			if err != nil {
				return err
			}
			p.children = append(p.children, cp)
			if err = cp.Init(); err != nil {
				return err
			}
			continue
		}
		// TODO: check on another type
		vl := newValue(v, p.rtype.Field(i))
		vl.owner = p
		p.values = append(p.values, vl)
	}
	return nil
}

func (p *parser) Parse() error {
	for _, v := range p.values {
		if err := v.define(); err != nil {
			return err
		}
	}
	for _, v := range p.children {
		if err := v.Parse(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) Path() string {
	if p.parent == nil {
		return ""
	}
	path := p.parent.Path()
	if path != "" {
		path += string(Separator)
	}
	return path + p.rtype.Name()
}
