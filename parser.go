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

func (e *EnvConf) fieldInitialized(f field) {
	cf := asConfigField(f)
	if cf == nil {
		return
	}
	dv, _ := cf.configuration.defaultValue.Value()
	e.opts.OnFieldInitialized(option.FieldInitializedArg{
		Name:         cf.name(),
		FullName:     cf.fullName(),
		Type:         cf.StructField.Type,
		Required:     cf.property.required,
		Description:  cf.property.description,
		FlagName:     cf.configuration.flag.Name(),
		EnvName:      cf.configuration.env.Name(),
		DefaultValue: dv,
	})
}

func (e *EnvConf) fieldDefined(f field) {
	cf := asConfigField(f)
	if cf == nil || !cf.isSet() {
		return
	}
	dv, _ := cf.configuration.defaultValue.Value()
	e.opts.OnFieldDefined(option.FieldDefinedArg{
		Name:         cf.name(),
		FullName:     cf.fullName(),
		Type:         cf.StructField.Type,
		Required:     cf.property.required,
		Description:  cf.property.description,
		FlagName:     cf.configuration.flag.Name(),
		EnvName:      cf.configuration.env.Name(),
		DefaultValue: dv,
		Value:        cf.value,
		Source:       cf.source,
	})
}

func (e *EnvConf) fieldNotDefined(f field, err error) {
	cf := asConfigField(f)
	if cf == nil {
		return
	}
	e.opts.OnFieldDefineErr(option.FieldDefineErrorArg{
		Name:     cf.name(),
		FullName: cf.fullName(),
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
