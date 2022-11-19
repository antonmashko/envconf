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
		Logger:   debugLogger,
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

func (e *EnvConf) fieldDefined(f field) {
	pt, ok := f.(*primitiveType)
	if !ok || pt.definedValue == nil {
		return
	}
	debugLogger.Printf("envconf: set variable name=%s value=%v source=%s",
		fullname(pt), pt.definedValue.value, pt.definedValue.source)
}

func (e *EnvConf) fieldNotDefined(f field, err error) {
	e.Logger.Printf("skipping error due not required field. field=%s err=%s",
		fullname(f.(Value)), err)
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
