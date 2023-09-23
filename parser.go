package envconf

import (
	"errors"
	"flag"

	"github.com/antonmashko/envconf/option"
)

// ErrNilData mean that exists nil pointer inside data struct
var ErrNilData = errors.New("nil data")

type EnvConf struct {
	external *externalConfig
	opts     *option.Options
}

func New() *EnvConf {
	return NewWithExternal(emptyExt{})
}

func NewWithExternal(e External) *EnvConf {
	return &EnvConf{
		external: newExternalConfig(e),
		opts:     &option.Options{},
	}
}

func (e *EnvConf) fieldInitialized(f field) {
	pt, ok := f.(*primitiveType)
	if !ok {
		return
	}
	dv, _ := pt.def.Value()
	e.opts.OnFieldInitialized(option.FieldInitializedArg{
		Name:         pt.name(),
		FullName:     fullname(pt),
		Type:         pt.sf.Type,
		Required:     pt.required,
		Description:  pt.desc,
		FlagName:     pt.flag.Name(),
		EnvName:      pt.env.Name(),
		DefaultValue: dv,
	})
}

func (e *EnvConf) fieldDefined(f field) {
	pt, ok := f.(*primitiveType)
	if !ok {
		return
	}
	if pt.definedValue == nil {
		return
	}
	dv, _ := pt.def.Value()
	e.opts.OnFieldDefined(option.FieldDefinedArg{
		Name:         pt.name(),
		FullName:     fullname(pt),
		Type:         pt.sf.Type,
		Required:     pt.required,
		Description:  pt.desc,
		FlagName:     pt.flag.Name(),
		EnvName:      pt.env.Name(),
		DefaultValue: dv,
		Source:       pt.definedValue.source,
		Value:        pt.definedValue.value,
	})
}

func (e *EnvConf) fieldNotDefined(f field, err error) {
	pt, ok := f.(*primitiveType)
	if !ok {
		return
	}
	e.opts.OnFieldDefineErr(option.FieldDefineErrorArg{
		Name:     pt.name(),
		FullName: fullname(pt),
		Err:      err,
	})
}

// Parse define variables inside data from different sources,
// such as flag/environment variable or default value
func (e *EnvConf) Parse(data interface{}, opts ...option.ClientOption) error {
	if data == nil {
		return ErrNilData
	}
	// enable help output
	option.WithCustomUsage().Apply(e.opts)
	for i := range opts {
		opts[i].Apply(e.opts)
	}
	p, err := newParentStructType(data, e)
	if err != nil {
		return err
	}
	e.external.setParentStruct(p)
	if err = p.init(); err != nil {
		return err
	}
	if e.opts.Usage() != nil {
		flag.Usage = e.opts.Usage()
	}
	flag.Parse()
	if fp := e.opts.FlagParsed(); fp != nil {
		if err = fp(); err != nil {
			return err
		}
	}
	if err = e.external.Unmarshal(data); err != nil {
		return err
	}
	return p.define()
}

// PriorityOrder return parsing priority order
func (e *EnvConf) PriorityOrder() []option.ConfigSource {
	return e.opts.PriorityOrder()
}

// Parse define variables inside data from different sources,
// such as flag/environment variable or default value
func Parse(data interface{}, opts ...option.ClientOption) error {
	return New().Parse(data, opts...)
}

// ParseWithExternal works same as Parse method but also can be used external sources
// (config files, key-value storages, etc.).
func ParseWithExternal(data interface{}, external External, opts ...option.ClientOption) error {
	return NewWithExternal(external).Parse(data, opts...)
}
