package option

import (
	"reflect"
	"testing"
)

func TestOptions_Callbacks_Ok(t *testing.T) {
	var (
		fi  bool
		fd  bool
		fde bool
	)
	opts := &Options{
		onFieldInitialized: func(arg FieldInitializedArg) { fi = true },
		onFieldDefined:     func(arg FieldDefinedArg) { fd = true },
		onFieldDefineErr:   func(arg FieldDefineErrorArg) { fde = true },
	}
	opts.OnFieldInitialized(FieldInitializedArg{})
	if !fi {
		t.Fatal("onFieldInitialized: not invoked")
	}
	opts.OnFieldDefined(FieldDefinedArg{})
	if !fd {
		t.Fatal("onFieldDefined: not invoked")
	}
	opts.OnFieldDefineErr(FieldDefineErrorArg{})
	if !fde {
		t.Fatal("onFieldDefineErr: not invoked")
	}
}

func TestOptions_DefaultPriorityOrder_Ok(t *testing.T) {
	opts := &Options{}
	if !reflect.DeepEqual(opts.PriorityOrder(), []ConfigSource{FlagVariable, EnvVariable, ExternalSource, DefaultValue}) {
		t.Fatal("unexpected result:", opts.PriorityOrder())
	}
}
