package envconf

import (
	"flag"
	"os"
	"reflect"
	"strings"
)

type configSource interface {
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

func (s *defaultValueSource) Value() (string, bool) {
	return s.v, s.defined
}
