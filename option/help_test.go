package option

import (
	"bytes"
	"flag"
	"reflect"
	"testing"
)

func TestWithCustomUsage_Ok(t *testing.T) {
	opt := WithCustomUsage()
	buff := bytes.NewBuffer([]byte{})
	flag.CommandLine.SetOutput(buff)
	opts := &Options{}
	opt.Apply(opts)
	if opts.onFieldInitialized == nil {
		t.Fatal("opts.onFieldInitialized is nil")
	}

	opts.OnFieldInitialized(FieldInitializedArg{
		Name:         "foo",
		FullName:     "foo",
		Type:         reflect.TypeOf(""),
		Required:     true,
		Description:  "desc",
		FlagName:     "-foo",
		EnvName:      "ENV_FOO",
		DefaultValue: "bar",
	})

	opts.Usage()()

	if buff.String() != "" {
		t.Fatal("unexpected result: ", buff.String())
	}
}

func TestWithoutCustomUsage_Ok(t *testing.T) {
	opt := WithCustomUsage()
	opts := &Options{}
	opt.Apply(opts)
	if opts.onFieldInitialized == nil || opts.Usage() == nil {
		t.Fatal("opts.onFieldInitialized or opts.Usage() is nil")
	}
	opt = WithoutCustomUsage()
	opt.Apply(opts)
	if opts.onFieldInitialized != nil || opts.Usage() != nil {
		t.Fatal("opts.onFieldInitialized or opts.Usage() is not nil")
	}
}
