package option

import (
	"reflect"
	"testing"
)

func TestConfigSourceString_Ok(t *testing.T) {
	if FlagVariable.String() != "Flag" {
		t.Fatal("[flag] unexpected result: ", FlagVariable.String())
	}
	if EnvVariable.String() != "Environment" {
		t.Fatal("[env] unexpected result: ", EnvVariable.String())
	}
	if ExternalSource.String() != "External" {
		t.Fatal("[external] unexpected result: ", ExternalSource.String())
	}
	if DefaultValue.String() != "Default" {
		t.Fatal("[default] unexpected result: ", DefaultValue.String())
	}
	if ConfigSource(123123).String() != "" {
		t.Fatal("[invalid] unexpected result: ", ConfigSource(123123).String())
	}
}

func TestWithPriorityOrder_Ok(t *testing.T) {
	opt := WithPriorityOrder(FlagVariable, EnvVariable, DefaultValue)
	if opt == nil {
		t.Fatal("opt is nil")
	}
	opts := &Options{}
	opt.Apply(opts)
	if len(opts.PriorityOrder()) != 3 {
		t.Fatal("unexpected result. len: ", len(opts.PriorityOrder()))
	}
}

func TestWithPriorityOrder_CustomOrder_Ok(t *testing.T) {
	opt := WithPriorityOrder(ExternalSource, EnvVariable, FlagVariable, DefaultValue)
	if opt == nil {
		t.Fatal("opt is nil")
	}
	opts := &Options{}
	opt.Apply(opts)
	if !reflect.DeepEqual(opts.PriorityOrder(), []ConfigSource{ExternalSource, EnvVariable, FlagVariable, DefaultValue}) {
		t.Fatal("unexpected result: ", opts.PriorityOrder())
	}
}

func TestWithPriorityOrder_Duplicates_Ok(t *testing.T) {
	opt := WithPriorityOrder(EnvVariable, EnvVariable, DefaultValue, EnvVariable, DefaultValue)
	if opt == nil {
		t.Fatal("opt is nil")
	}
	opts := &Options{}
	opt.Apply(opts)
	if !reflect.DeepEqual(opts.PriorityOrder(), []ConfigSource{EnvVariable, DefaultValue}) {
		t.Fatal("unexpected result: ", opts.PriorityOrder())
	}
}

func TestWithPriorityOrder_DefaultOrder_Ok(t *testing.T) {
	opt := WithPriorityOrder()
	if opt == nil {
		t.Fatal("opt is nil")
	}
	opts := &Options{}
	opt.Apply(opts)
	if !reflect.DeepEqual(opts.PriorityOrder(), []ConfigSource{FlagVariable, EnvVariable, ExternalSource, DefaultValue}) {
		t.Fatal("unexpected result: ", opts.PriorityOrder())
	}
}

func TestWithPriorityOrder_InvalidSource_Ok(t *testing.T) {
	opt := WithPriorityOrder(ConfigSource(123), ConfigSource(124))
	if opt == nil {
		t.Fatal("opt is nil")
	}
	opts := &Options{}
	opt.Apply(opts)
	if !reflect.DeepEqual(opts.PriorityOrder(), []ConfigSource{FlagVariable, EnvVariable, ExternalSource, DefaultValue}) {
		t.Fatal("unexpected result: ", opts.PriorityOrder())
	}
}

func TestWithPriorityOrder_InvalidSourceAndValid_Ok(t *testing.T) {
	opt := WithPriorityOrder(ConfigSource(123), ExternalSource)
	if opt == nil {
		t.Fatal("opt is nil")
	}
	opts := &Options{}
	opt.Apply(opts)
	if !reflect.DeepEqual(opts.PriorityOrder(), []ConfigSource{ExternalSource}) {
		t.Fatal("unexpected result: ", opts.PriorityOrder())
	}
}
