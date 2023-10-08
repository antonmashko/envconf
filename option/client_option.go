package option

import (
	"reflect"

	"github.com/antonmashko/envconf/external"
)

type FieldInitializedArg struct {
	Name        string
	FullName    string
	Type        reflect.Type
	Required    bool
	Description string

	FlagName     string
	EnvName      string
	DefaultValue interface{}
}

type FieldDefinedArg struct {
	Name        string
	FullName    string
	Type        reflect.Type
	Required    bool
	Description string

	FlagName     string
	EnvName      string
	DefaultValue interface{}

	Source ConfigSource
	Value  interface{}
}

type FieldDefineErrorArg struct {
	Name     string
	FullName string
	Type     reflect.Type
	Err      error
}

type ClientOption interface {
	Apply(*Options)
}

type Options struct {
	external           external.External
	priorityOrder      []ConfigSource
	flagParsed         func() error
	usage              func()
	onFieldInitialized func(FieldInitializedArg)
	onFieldDefined     func(FieldDefinedArg)
	onFieldDefineErr   func(FieldDefineErrorArg)
	externalInjection  func(string) (string, ConfigSource)
}

func (o *Options) External() external.External {
	return o.external
}

func (o *Options) PriorityOrder() []ConfigSource {
	if len(o.priorityOrder) == 0 {
		return []ConfigSource{
			FlagVariable,
			EnvVariable,
			ExternalSource,
			DefaultValue,
		}
	}
	return o.priorityOrder
}

func (o *Options) Usage() func() {
	return o.usage
}

func (o *Options) FlagParsed() func() error {
	return o.flagParsed
}

func (o *Options) OnFieldInitialized(arg FieldInitializedArg) {
	if o.onFieldInitialized != nil {
		o.onFieldInitialized(arg)
	}
}

func (o *Options) OnFieldDefined(arg FieldDefinedArg) {
	if o.onFieldDefined != nil {
		o.onFieldDefined(arg)
	}
}

func (o *Options) OnFieldDefineErr(arg FieldDefineErrorArg) {
	if o.onFieldDefineErr != nil {
		o.onFieldDefineErr(arg)
	}
}

func (o *Options) ExternalInjection() func(string) (string, ConfigSource) {
	return o.externalInjection
}
