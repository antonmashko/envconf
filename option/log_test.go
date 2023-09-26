package option

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
)

type testPrinter struct {
	message string
}

func (p *testPrinter) Print(arg ...interface{}) {
	p.message = fmt.Sprint(arg...)
}

func TestWithLog_Ok(t *testing.T) {
	opt := WithLog(log.Default())
	if opt == nil {
		t.Fatal("nil option")
	}
	opts := &Options{}
	opt.Apply(opts)
	if opts.onFieldDefined == nil {
		t.Fatal("opts.onFieldDefined is nil")
	}
	if opts.onFieldDefineErr == nil {
		t.Fatal("opts.onFieldDefineErr is nil")
	}
}

func TestWithLog_Nil_Ok(t *testing.T) {
	opt := WithLog(nil)
	if opt == nil {
		t.Fatal("nil option")
	}
	opt.Apply(&Options{})
}

func TestWithLog_PrintMessage_Ok(t *testing.T) {
	p := &testPrinter{}
	opt := WithLog(p)
	if opt == nil {
		t.Fatal("nil option")
	}
	opt.printDefined(FieldDefinedArg{
		Name:     "foo",
		FullName: "foo",
		Type:     reflect.TypeOf("foo"),
		Source:   EnvVariable,
		Value:    "bar",
	})

	if p.message != `field="foo" value="bar" type="string" source="Environment"` {
		t.Fatal("unexpected result: ", p.message)
	}
}

func TestWithLog_PrintWithSecret_Ok(t *testing.T) {
	p := &testPrinter{}
	opt := WithLog(p)
	if opt == nil {
		t.Fatal("nil option")
	}
	opt.printDefined(FieldDefinedArg{
		Name:     "Password",
		FullName: "User.Password",
		Type:     reflect.TypeOf(""),
		Source:   EnvVariable,
		Value:    "secret_password",
	})

	if p.message != `field="User.Password" value="******" type="string" source="Environment"` {
		t.Fatal("unexpected result: ", p.message)
	}
}

func TestWithLog_PrintErr_Ok(t *testing.T) {
	p := &testPrinter{}
	opt := WithLog(p)
	if opt == nil {
		t.Fatal("nil option")
	}
	opt.printErr(FieldDefineErrorArg{
		FullName: "foo",
		Err:      errors.New("custom error"),
	})

	if p.message != `field="foo" err="custom error"` {
		t.Fatal("unexpected result: ", p.message)
	}
}
