package envconf

import (
	"errors"
	"flag"
)

var defaultEnvConf *EnvConf

// ErrNilData mean that exists nil pointer inside data struct
var ErrNilData = errors.New("nil data")

// FlagParsed define this callback when you need handle flags
// This callback will raise after method flag.Parse()
// return not nil error interrupt pasring
var FlagParsed func() error

func init() {
	defaultEnvConf = New()
}

type EnvConf struct {
	Logger   Logger
	external *externalConfig
	priority []ConfigSource
	help     *help
}

func New() *EnvConf {
	return NewWithExternal(emptyExt{})
}

func NewWithExternal(e External) *EnvConf {
	var h *help
	if UseCustomHelp {
		h = &help{
			fields: make([]*primitiveType, 0),
		}
	}

	return &EnvConf{
		external: newExternalConfig(e),
		help:     h,
		Logger:   debugLogger,
		priority: []ConfigSource{
			FlagVariable, EnvVariable, ExternalSource, DefaultValue,
		},
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
	e.Logger.Printf("envconf: set variable name=%s value=%v source=%s",
		fullname(pt), pt.definedValue.value, pt.definedValue.source)
}

func (e *EnvConf) fieldNotDefined(f field, err error) {
	e.Logger.Printf("skipping error due not required field. field=%s err=%s",
		fullname(f), err)
}

// Parse define variables inside data from different sources,
// such as flag/environment variable or default value
func (e *EnvConf) Parse(data interface{}) error {
	if data == nil {
		return ErrNilData
	}
	p, err := newParentStructType(data, e)
	if err != nil {
		return err
	}
	e.external.setParentStruct(p)
	if err = p.init(); err != nil {
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
	return p.define()
}

// SetPriorityOrder overrides default priority order.
// Default priority order is: Flag, Environment variable, External source, Default value.
func (e *EnvConf) SetPriorityOrder(s ...ConfigSource) {
	if len(s) == 0 {
		return
	}
	po := make(map[ConfigSource]int)
	var idx int
	for _, p := range s {
		if p != FlagVariable && p != EnvVariable &&
			p != ExternalSource && p != DefaultValue {
			continue
		}
		if _, ok := po[p]; !ok {
			po[p] = idx
			idx++
		}
	}
	if len(po) == 0 {
		return
	}
	result := make([]ConfigSource, len(po))
	for s, idx := range po {
		result[idx] = s
	}
	e.priority = result
}

// PriorityOrder return parsing priority order
func (e *EnvConf) PriorityOrder() []ConfigSource {
	return e.priority
}

// Parse define variables inside data from different sources,
// such as flag/environment variable or default value
func Parse(data interface{}) error {
	return defaultEnvConf.Parse(data)
}

// ParseWithExternal works same as Parse method but also can be used external sources
// (config files, key-value storages, etc.).
func ParseWithExternal(data interface{}, external External) error {
	ecfg := NewWithExternal(external)
	ecfg.Logger = defaultEnvConf.Logger
	ecfg.SetPriorityOrder(defaultEnvConf.PriorityOrder()...)
	return ecfg.Parse(data)
}

// SetLogger define debug logger.
// This logger will print defined values in data fields
func SetLogger(logger Logger) {
	if logger != nil {
		defaultEnvConf.Logger = logger
	}
}

// SetPriority overrides default priority order.
// Default priority order is: Flag, Environment variable, External source, Default value.
func SetPriority(s ...ConfigSource) {
	defaultEnvConf.SetPriorityOrder(s...)
}
