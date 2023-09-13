package envconf

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"strings"
)

type ConfigSource int

const (
	FlagVariable ConfigSource = iota
	EnvVariable
	ExternalSource
	DefaultValue
)

func (s ConfigSource) String() string {
	switch s {
	case FlagVariable:
		return "Flag"
	case EnvVariable:
		return "Environment"
	case ExternalSource:
		return "External"
	case DefaultValue:
		return "Default"
	}
	return ""
}

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

var (
	//errors
	errInvalidFiled              = errors.New("invalid field")
	errFiledIsNotSettable        = errors.New("field is not settable")
	ErrUnsupportedType           = errors.New("unsupported type")
	errConfigurationNotSpecified = errors.New("configuration not specified")
)

// Var is configuration variable for defining primitive data types
type Var interface {
	Name() string
	Value() (string, bool)
}

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

func (s *flagSource) Value() (string, bool) {
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

func (s *envSource) Name() string {
	return s.name
}

func (s *envSource) Value() (string, bool) {
	if s.name != tagIgnored {
		return os.LookupEnv(s.name)
	}
	return "", false
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

func (s *defaultValueSource) Name() string {
	return tagDefault
}

func (s *defaultValueSource) Value() (string, bool) {
	return s.v, s.defined
}

type configVar struct {
	flag Var
	env  Var
	ext  Var
	def  Var
}

func newConfigVar(f field, tag reflect.StructField, extConf *externalConfig) *configVar {
	return &configVar{}
}
