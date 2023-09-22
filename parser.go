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
		opts: &option.Options{
			PriorityOrder: []option.ConfigSource{
				option.FlagVariable, option.EnvVariable, option.ExternalSource, option.DefaultValue,
			},
		},
	}
}

func (e *EnvConf) fieldInitialized(f field) {
	if e.opts.OnFieldInitialized == nil || len(e.opts.OnFieldInitialized) == 0 {
		return
	}
	pt, ok := f.(*primitiveType)
	if !ok {
		return
	}
	dv, _ := pt.def.Value()
	for i := range e.opts.OnFieldInitialized {
		e.opts.OnFieldInitialized[i](option.FieldInitializedArg{
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
}

func (e *EnvConf) fieldDefined(f field) {
	if e.opts.OnFieldDefined == nil || len(e.opts.OnFieldDefined) == 0 {
		return
	}
	pt, ok := f.(*primitiveType)
	if !ok {
		return
	}
	if pt.definedValue == nil {
		return
	}
	dv, _ := pt.def.Value()
	for i := range e.opts.OnFieldDefined {
		e.opts.OnFieldDefined[i](option.FieldDefinedArg{
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
}

func (e *EnvConf) fieldNotDefined(f field, err error) {
	// e.printErr(f, err, "skipping error because the field is not required.")
	if e.opts.OnFieldDefineErr == nil || len(e.opts.OnFieldDefineErr) == 0 {
		return
	}
	pt, ok := f.(*primitiveType)
	if !ok {
		return
	}

	for i := range e.opts.OnFieldDefineErr {
		e.opts.OnFieldDefineErr[i](option.FieldDefineErrorArg{
			Name:     pt.name(),
			FullName: fullname(pt),
			Err:      err,
		})
	}
}

// Parse define variables inside data from different sources,
// such as flag/environment variable or default value
func (e *EnvConf) Parse(data interface{}, opts ...option.ClientOption) error {
	if data == nil {
		return ErrNilData
	}
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
	if e.opts.Usage != nil {
		flag.Usage = e.opts.Usage
	}
	flag.Parse()
	if e.opts.FlagParsed != nil {
		if err = e.opts.FlagParsed(); err != nil {
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
	return e.opts.PriorityOrder
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

// SetLogger define debug printer.
// This logger will print defined values in data fields
// func SetOptions(opts ...option.ClientOption) {
// 	for i := range opts {
// 		if opts[i] != nil {
// 			opts[i].Apply(defaultEnvConf.opts)
// 		}
// 	}
// }
