package envconf

import (
	"flag"
	"os"
	"reflect"
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

func newFlagSource(f field, tag reflect.StructField, usage string) *flagSource {
	name, ok := tag.Tag.Lookup(tagFlag)
	if !ok || name == tagNotDefined {
		name = tagIgnored
	} else if name == valDefault {
		// generating flag name
		const flagDelim = "-"
		name = strings.ToLower(strings.ReplaceAll(fullname(f), fieldNameDelim, flagDelim))
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

func (s *flagSource) Value() (interface{}, bool) {
	if s.name == tagIgnored {
		return "", false
	}
	return s.v, s.defined
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

func newEnvSource(f field, tag reflect.StructField) *envSource {
	name, ok := tag.Tag.Lookup(tagEnv)
	if !ok || name == tagNotDefined {
		name = tagIgnored
	} else if name == valDefault {
		// generating env var name
		const envDelim = "_"
		name = strings.ToUpper(strings.ReplaceAll(fullname(f), fieldNameDelim, envDelim))
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

func (s *envSource) Value() (interface{}, bool) {
	if s.name != tagIgnored {
		return os.LookupEnv(s.name)
	}
	return "", false
}

type externalValueSource struct {
	f                 field
	allowEnvInjection bool
}

func newExternalValueSource(f field, allowEnvInjection bool) *externalValueSource {
	return &externalValueSource{
		f:                 f,
		allowEnvInjection: allowEnvInjection,
	}
}

func (s *externalValueSource) Value() (interface{}, option.ConfigSource, bool) {
	if s.f.parent() == nil {
		return nil, option.ExternalSource, false
	}
	es := s.f.parent().externalSource()
	v, ok := es.Read(s.f.structField().Name)
	if !ok {
		return nil, option.ExternalSource, false
	}
	if !s.allowEnvInjection {
		return v, option.ExternalSource, true
	}
	str, ok := v.(string)
	if !ok {
		return v, option.ExternalSource, true
	}
	envInj := newEnvInjection(str)
	if envInj == nil {
		return v, option.ExternalSource, true
	}
	var ev interface{}
	ev, ok = envInj.Value()
	if !ok {
		return v, option.ExternalSource, true
	}
	return ev, option.EnvVariable, true
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

func (s *defaultValueSource) Value() (interface{}, bool) {
	return s.v, s.defined
}
