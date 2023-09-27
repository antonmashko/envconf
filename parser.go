package envconf

import (
	"flag"

	"github.com/antonmashko/envconf/external"
	"github.com/antonmashko/envconf/option"
)

type EnvConf struct {
	opts *option.Options
}

func New() *EnvConf {
	return &EnvConf{
		opts: &option.Options{},
	}
}

func (e *EnvConf) fieldType(f field) *fieldType {
	switch ft := f.(type) {
	case *fieldType:
		return ft
	case *ptrType:
		return e.fieldType(ft.field)
	case *interfaceType:
		return e.fieldType(ft.field)
	default:
		return nil
	}
}

func (e *EnvConf) fieldInitialized(f field) {
	pt := e.fieldType(f)
	if pt == nil {
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
	pt := e.fieldType(f)
	if pt == nil {
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
	pt := e.fieldType(f)
	if pt == nil {
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

	extMapper := external.NewExternalConfigMapper(e.opts.External())
	p, err := newParentStructType(data, e)
	if err != nil {
		return err
	}
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
	if err = extMapper.Unmarshal(data); err != nil {
		return err
	}
	p.ext = extMapper.Data()
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
