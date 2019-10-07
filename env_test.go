package envconf

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestSimpleExternalEnvConfigOK(t *testing.T) {
	env := bytes.NewBuffer([]byte("Foo=Bar"))
	tc := struct {
		Foo string `env:"Foo"`
	}{}
	NewEnvConf().SetEnv(env)
	err := Parse(&tc)
	if err != nil {
		t.Error(err)
	}
	if tc.Foo != "Bar" {
		t.Errorf("incorrect values was set. %#v", tc.Foo)
	}
	fmt.Println(os.Getenv("Name"))
}

func TestSimpleExternalEnvConfigFieldWithUnderscoreOK(t *testing.T) {
	env := bytes.NewBuffer([]byte("FOO_BAR=FOO_BAR"))
	tc := struct {
		FooBar string `env:"FOO_BAR"`
	}{}
	NewEnvConf().SetEnv(env)
	err := Parse(&tc)
	if err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.FooBar != "FOO_BAR" {
		t.Errorf("incorrect values was set. %#v", tc.FooBar)
	}
}

func TestSimpleExternaEnvConfigTwoVariablesOK(t *testing.T) {
	env := bytes.NewBuffer([]byte("FOO_FIRST=FOO_FIRST\nFOO_SECOND=FOO_SECOND\n"))
	tc := struct {
		FooFirst  string `env:"FOO_FIRST"`
		FooSecond string `env:"FOO_SECOND"`
	}{}
	NewEnvConf().SetEnv(env)
	err := Parse(&tc)
	if err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.FooFirst != "FOO_FIRST" {
		t.Errorf("incorrect first values was set. %#v", tc.FooFirst)
	}
	if tc.FooSecond != "FOO_SECOND" {
		t.Errorf("incorrect value was set. %#v", tc.FooSecond)
	}
}
