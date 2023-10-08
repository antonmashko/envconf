package envconf

import (
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/antonmashko/envconf/option"
)

const (
	tagFlag        = "flag"
	tagEnv         = "env"
	tagDefault     = "default"
	tagRequired    = "required"
	tagDescription = "description"
	tagIgnored     = "-"
	tagNotDefined  = ""

	valIgnored    = "ignored"
	valNotDefined = "N/D"
	valDefault    = "*"
)

type flagSource struct {
	name    string
	v       string
	defined bool
}

func newFlagSource(f *configField, tag reflect.StructField, usage string) *flagSource {
	name, ok := tag.Tag.Lookup(tagFlag)
	if !ok || name == tagNotDefined {
		name = tagIgnored
	} else if name == valDefault {
		// generating flag name
		const flagDelim = "-"
		name = strings.ToLower(fullname(f, flagDelim))
	}
	fs := &flagSource{
		name: name,
	}
	if name != tagIgnored {
		flag.Var(fs, name, usage)
	}
	return fs
}

func (s *flagSource) Name() string {
	return s.name
}

func (s *flagSource) Value() (interface{}, option.ConfigSource) {
	if s.name == tagIgnored {
		return "", option.NoConfigValue
	}
	if !s.defined {
		return "", option.NoConfigValue
	}
	return s.v, option.FlagVariable
}

func (s *flagSource) Set(value string) error {
	s.v = value
	s.defined = true
	return nil
}

func (s *flagSource) String() string {
	return s.v
}

type envSource struct {
	name string
}

func newEnvSource(f *configField, tag reflect.StructField) *envSource {
	name, ok := tag.Tag.Lookup(tagEnv)
	if !ok || name == tagNotDefined {
		name = tagIgnored
	} else if name == valDefault {
		// generating env var name
		const envDelim = "_"
		name = strings.ToUpper(fullname(f, envDelim))
	}
	return &envSource{
		name: name,
	}
}

func newEnvInjection(value string) *envSource {
	var ok bool
	const prefix = "${"
	value, ok = strings.CutPrefix(value, prefix)
	if !ok {
		return nil
	}

	const suffix = "}"
	value, ok = strings.CutSuffix(value, suffix)
	if !ok {
		return nil
	}

	value = strings.TrimSpace(value)
	const envInjectionPrefix = ".env."
	value, ok = strings.CutPrefix(value, envInjectionPrefix)
	if !ok {
		return nil
	}
	return &envSource{name: value}
}

func (s *envSource) Name() string {
	return s.name
}

func (s *envSource) Value() (interface{}, option.ConfigSource) {
	if s.name == tagIgnored {
		return "", option.NoConfigValue
	}
	v, ok := os.LookupEnv(s.name)
	if !ok {
		return "", option.NoConfigValue
	}
	return v, option.EnvVariable
}

type externalSource struct {
	f                 field
	allowEnvInjection bool
}

func newExternalSource(f field, allowEnvInjection bool) *externalSource {
	return &externalSource{
		f:                 f,
		allowEnvInjection: allowEnvInjection,
	}
}

func (s *externalSource) Value() (interface{}, option.ConfigSource) {
	if s.f.parent() == nil {
		return nil, option.NoConfigValue
	}
	v, ok := s.f.parent().externalSource().
		Read(s.f.structField().Name)
	if !ok {
		return nil, option.NoConfigValue
	}
	if !s.allowEnvInjection {
		return v, option.ExternalSource
	}
	str, ok := v.(string)
	if !ok {
		return v, option.ExternalSource
	}
	envInj := newEnvInjection(str)
	if envInj == nil {
		return v, option.ExternalSource
	}
	return envInj.Value()
}

type defaultValueSource struct {
	defined bool
	v       string
}

func newDefaultValueSource(tag reflect.StructField) *defaultValueSource {
	var s defaultValueSource
	s.v, s.defined = tag.Tag.Lookup(tagDefault)
	return &s
}

func (s *defaultValueSource) Value() (interface{}, option.ConfigSource) {
	if !s.defined {
		return nil, option.NoConfigValue
	}
	return s.v, option.DefaultValue
}

type configField struct {
	reflect.StructField
	parentField   field
	parser        *EnvConf
	configuration struct {
		flag         *flagSource
		env          *envSource
		external     *externalSource
		defaultValue *defaultValueSource
	}
	property struct {
		required    bool
		description string
	}
	value  interface{}
	source option.ConfigSource
}

func newConfigField(parent field, sf reflect.StructField, parser *EnvConf) *configField {
	return newDefinedConfigField(nil, option.NoConfigValue, parent, sf, parser)
}

func newDefinedConfigField(v interface{}, cs option.ConfigSource, parent field, sf reflect.StructField, parser *EnvConf) *configField {
	return &configField{
		StructField: sf,
		parentField: parent,
		parser:      parser,
		value:       v,
		source:      cs,
	}
}

// initialize setting for specific field
func (f *configField) init(fl field) error {
	req, ok := f.Tag.Lookup(tagRequired)
	if ok {
		var err error
		f.property.required, err = strconv.ParseBool(req)
		if err != nil {
			return err
		}
	}
	f.property.description = f.Tag.Get(tagDescription)
	f.configuration.flag = newFlagSource(f, f.StructField, f.property.description)
	f.configuration.env = newEnvSource(f, f.StructField)
	f.configuration.external = newExternalSource(fl, f.parser.opts.AllowExternalEnvInjection)
	f.configuration.defaultValue = newDefaultValueSource(f.StructField)
	return nil
}

func (f *configField) name() string {
	return f.Name
}

func (f *configField) fullName() string {
	return fullname(f, fieldNameDelim)
}

func (f *configField) parent() field {
	if f.parentField == nil {
		return nil
	}
	return f.parentField
}

func (f *configField) structField() reflect.StructField {
	return f.StructField
}

func (f *configField) IsRequired() bool {
	return f.property.required
}

func (f *configField) isSet() bool {
	return f.value != nil && f.source != option.NoConfigValue
}

func (f *configField) set(v interface{}, cs option.ConfigSource) error {
	f.value = v
	f.source = cs
	return nil
}

func (f *configField) Value() (interface{}, option.ConfigSource) {
	if f.isSet() {
		return f.value, f.source
	}
	priority := f.parser.PriorityOrder()
	for _, p := range priority {
		var confF func() (interface{}, option.ConfigSource) = nil
		switch p {
		case option.FlagVariable:
			confF = f.configuration.flag.Value
		case option.EnvVariable:
			confF = f.configuration.env.Value
		case option.ExternalSource:
			confF = f.configuration.external.Value
		case option.DefaultValue:
			confF = f.configuration.defaultValue.Value
		}
		if confF == nil {
			continue
		}
		v, cs := confF()
		if cs != option.NoConfigValue {
			return v, cs
		}
	}
	return nil, option.NoConfigValue
}
