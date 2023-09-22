package option

import "reflect"

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

type Options struct {
	PriorityOrder      []ConfigSource
	FlagParsed         func() error
	Usage              func()
	OnFieldInitialized []func(FieldInitializedArg)
	OnFieldDefined     []func(FieldDefinedArg)
	OnFieldDefineErr   []func(FieldDefineErrorArg)
}

type ClientOption interface {
	Apply(*Options)
}
