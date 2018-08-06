package envconf

import (
	"errors"
	"flag"
	"log"
	"os"
	"reflect"
)

const Separator = "."

// Errors
var (
	ErrNilData = errors.New("nil data")
)

// OnFlagParsed define this callback when you need handle flags
// This callback will raise after method flag.Parse()
// return not nil error interupt pasring
var OnFlagParsed func() error

// Loggers
var (
	traceLogger Logger = &logger{l: log.New(os.Stdout, "", log.Ltime)}
	errorLogger Logger = &logger{l: log.New(os.Stderr, "", log.Ltime)}
)

func SetTraceLogger(logger Logger) {
	if logger != nil {
		traceLogger = logger
	}
}

func SetErrorLogger(logger Logger) {
	if logger != nil {
		errorLogger = logger
	}
}

//Parse fiend tag(annotations) for each field as set value
func Parse(data interface{}) error {
	return ParseWithExternal(data, &emptyConfig{})
}

func ParseWithExternal(data interface{}, external External) error {
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
	if OnFlagParsed != nil {
		if err = OnFlagParsed(); err != nil {
			return err
		}
	}
	return p.Parse()
}

type parser struct {
	value    reflect.Value
	rtype    reflect.Type
	tag      reflect.StructField
	parent   *parser
	external External
	children []*parser
	values   []*value
}

func newParser(data interface{}, external External) (*parser, error) {
	return newChildParser(nil, reflect.ValueOf(data), external)
}

func newChildParser(parent *parser, rvalue reflect.Value, external External) (*parser, error) {
	p := &parser{external: external}
	// get value from pointer
	for rvalue.Kind() == reflect.Ptr {
		// check on nil
		if rvalue.IsNil() {
			return nil, ErrNilData
		}
		rvalue = rvalue.Elem()
	}
	p.value = rvalue
	p.rtype = p.value.Type() // remember type
	p.children = make([]*parser, 0)
	p.values = make([]*value, 0)
	p.parent = parent
	return p, nil
}

func (p *parser) Init() error {
	for i := 0; i < p.value.NumField(); i++ {
		v := p.value.Field(i)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return ErrNilData
			}
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
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
		vl := newValue(p, v, p.rtype.Field(i))
		vl.owner = p
		p.values = append(p.values, vl)
	}
	return nil
}

func (p *parser) Name() string {
	return p.tag.Name
}

func (p *parser) Tag() reflect.StructField {
	return p.tag
}

func (p *parser) Owner() Value {
	return p.parent
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
