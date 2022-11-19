package envconf

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
)

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

type EnvConf struct {
	Logger   Logger
	external External
	help     *help
}

func New() *EnvConf {
	return NewWithExternal(&emptyExt{})
}

func NewWithExternal(e External) *EnvConf {
	var h *help
	if UseCustomHelp {
		h = &help{
			fields: make([]*primitiveType, 0),
		}
	}

	return &EnvConf{
		external: e,
		help:     h,
	}
}

func (e *EnvConf) fieldInitialized(f field) {
	if e.help == nil {
		return
	}
	pt, ok := f.(*primitiveType)
	if !ok {
		return
	}
	e.help.fields = append(e.help.fields, pt)
}

func (e *EnvConf) Parse(data interface{}) error {
	if data == nil {
		return ErrNilData
	}
	p, err := newParentStructType(data, e)
	if err != nil {
		return err
	}
	if err = p.Init(); err != nil {
		return err
	}
	if e.help != nil {
		flag.Usage = e.help.usage
	}
	flag.Parse()
	if FlagParsed != nil {
		if err = FlagParsed(); err != nil {
			return err
		}
	}
	if err = e.external.Unmarshal(data); err != nil {
		return err
	}
	return p.Define()
}

// Parse define variables inside data from different sources,
// such as flag/environment variable or default value
func Parse(data interface{}) error {
	return New().Parse(data)
}

// ParseWithExternal works same as Parse method but also can be used external sources
// (config files, key-value storages, etc.).
func ParseWithExternal(data interface{}, external External) error {
	return NewWithExternal(external).Parse(data)
}

func fullname(v Value) string {
	name := v.Name()
	for {
		v = v.Owner()
		if v == nil {
			break
		}
		oname := v.Name()
		if oname != "" {
			name = oname + "." + name
		}
	}
	return name
}

// type parser struct {
// 	value    reflect.Value
// 	rtype    reflect.Type
// 	tag      reflect.StructField
// 	parent   *parser
// 	external External
// 	children []*parser
// 	values   []*value
// }

// func newParser(data interface{}, external External) (*parser, error) {
// 	v, err := depointerize(reflect.ValueOf(data))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return newChildParser(nil, v, reflect.StructField{}, external), err
// }

// func newChildParser(p *parser, v reflect.Value, tag reflect.StructField, e External) *parser {
// 	result := &parser{
// 		parent:   p,
// 		external: e,
// 		value:    v,
// 		rtype:    v.Type(),
// 		tag:      tag,
// 		children: make([]*parser, 0),
// 		values:   make([]*value, 0),
// 	}
// 	return result
// }

// func (p *parser) Init() error {
// 	for i := 0; i < p.value.NumField(); i++ {
// 		v, err := depointerize(p.value.Field(i))
// 		if err != nil {
// 			if IgnoreNilData && err == ErrNilData {
// 				continue
// 			}
// 			return err
// 		}
// 		tag := p.rtype.Field(i)
// 		if v.Kind() == reflect.Struct {
// 			cp := newChildParser(p, v, tag, p.external)
// 			p.children = append(p.children, cp)
// 			if err = cp.Init(); err != nil {
// 				return err
// 			}
// 			continue
// 		}
// 		// TODO: check on another type
// 		vl := newValue(p, v, tag)
// 		p.values = append(p.values, vl)
// 	}
// 	return nil
// }

// func (p *parser) Name() string {
// 	return p.tag.Name
// }

// func (p *parser) Tag() reflect.StructField {
// 	return p.tag
// }

// func (p *parser) Owner() Value {
// 	return p.parent
// }

// func (p *parser) Parse() error {
// 	for _, v := range p.values {
// 		if err := v.define(); err != nil {
// 			if v.required {
// 				return &Error{
// 					Message:   "failed to define field",
// 					Inner:     err,
// 					FieldName: v.fullname(),
// 				}
// 			}
// 			if err == errConfigurationNotSpecified {
// 				continue
// 			}
// 			debugLogger.Printf("skipping error due not required field. field=%s err=%s", v.fullname(), err)
// 		}
// 	}
// 	for _, v := range p.children {
// 		if err := v.Parse(); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
