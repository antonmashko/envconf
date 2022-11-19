package envconf

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"strings"
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

var (
	//errors
	errInvalidFiled              = errors.New("invalid field")
	errFiledIsNotSettable        = errors.New("field is not settable")
	errUnsupportedType           = errors.New("unsupported type")
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

func newFlagSource(tag reflect.StructField, usage string) *flagSource {
	name := tag.Tag.Get(tagFlag)
	if name == tagNotDefined {
		name = tagIgnored
	} else if strings.ToLower(name) == valDefault {
		name = strings.ToLower(tag.Name)
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

func newEnvSource(tag reflect.StructField) *envSource {
	name := tag.Tag.Get(tagEnv)
	if name == tagNotDefined {
		name = tagIgnored
	} else if name == valDefault {
		name = strings.ToUpper(tag.Name)
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
