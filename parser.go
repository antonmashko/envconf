package envconf

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

// IgnoreNilData throw ErrNilData error if pointer(s) inside data struct is nil
var IgnoreNilData = false

// ErrNilData mean that exists nil pointer inside data struct
var ErrNilData = errors.New("nil data")

// FlagParsed define this callback when you need handle flags
// This callback will raise after method flag.Parse()
// return not nil error interrupt pasring
var FlagParsed func() error

var debugLogger Logger = &logger{l: log.New(ioutil.Discard, "", log.Ltime)}

// SetLogger define debug logger.
// This logger will print setted values in data fields
func SetLogger(logger Logger) {
	if logger != nil {
		debugLogger = logger
	}
}

// Parse define variables inside data from different sources,
// such as flag/environment variable or default value
func Parse(data interface{}) error {
	return ParseWithExternal(data, &emptyExt{})
}

// ParseWithExternal works same as Parse method but also can be used external sources
// (config files, key-value storages, etc.).
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
	if FlagParsed != nil {
		if err = FlagParsed(); err != nil {
			return err
		}
	}
	if err = external.Unmarshal(data); err != nil {
		return err
	}
	return p.Parse()
}

func depointerize(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v
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
	v := depointerize(reflect.ValueOf(data))
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, ErrNilData
	}
	return newChildParser(nil, v, reflect.StructField{}, external), nil
}

func newChildParser(p *parser, v reflect.Value, tag reflect.StructField, e External) *parser {
	return &parser{
		parent:   p,
		external: e,
		value:    v,
		rtype:    v.Type(),
		tag:      tag,
		children: make([]*parser, 0),
		values:   make([]*value, 0),
	}
}

func (p *parser) Init() error {
	for i := 0; i < p.value.NumField(); i++ {
		var err error
		// in case on pointer
		var verifyOnSet bool
		var tmp reflect.Value

		v := depointerize(p.value.Field(i))
		tag := p.rtype.Field(i)

		vkind := v.Kind()
		// if vkind == reflect.Ptr {
		// 	tmp = v
		// 	v = reflect.New(v.Type().Elem()).Elem()
		// 	verifyOnSet = true
		// }

		if vkind == reflect.Struct {
			cp := newChildParser(p, v, tag, p.external)
			if err = cp.Init(); err != nil {
				return err
			}

			if verifyOnSet {
				if len(cp.children) <= 0 && len(cp.values) <= 0 {
					continue
				}
				tmp.Set(v) // updating value
			}
			p.children = append(p.children, cp)
			continue
		}

		vl := newValue(p, v, tag)
		if verifyOnSet {
			if v.IsZero() {
				continue
			}
			fmt.Printf("%s\n", v.Elem().Type())
			tmp.Set(v)
		}
		p.values = append(p.values, vl)
	}
	return nil
}

func (p *parser) Parse() error {
	for _, v := range p.values {
		if err := v.define(); err != nil {
			if v.required {
				return &Error{
					Message:   "failed to define field",
					Inner:     err,
					FieldName: v.fullname(),
				}
			}
			if err == errConfigurationNotSpecified {
				continue
			}
			debugLogger.Printf("skipping error due not required field. field=%s err=%s", v.fullname(), err)
		}
	}

	for _, v := range p.children {
		if err := v.Parse(); err != nil {
			return err
		}
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
