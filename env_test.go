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
	envf := NewEnvConf()
	if err := envf.Parse(env); err != nil {
		t.Error("invalid parse env data:", err)
	}
	if err := envf.Set(); err != nil {
		t.Error("invalid set env variables:", err)
	}
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
	envf := NewEnvConf()
	if err := envf.Parse(env); err != nil {
		t.Error("invalid parse env data:", err)
	}
	if err := envf.Set(); err != nil {
		t.Error("invalid set env variables:", err)
	}
	err := Parse(&tc)
	if err != nil {
		t.Error(err)
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
	envf := NewEnvConf()
	if err := envf.Parse(env); err != nil {
		t.Error("invalid parse env data:", err)
	}
	if err := envf.Set(); err != nil {
		t.Error("invalid set env variables:", err)
	}
	err := Parse(&tc)
	if err != nil {
		t.Error(err)
	}
	if tc.FooFirst != "FOO_FIRST" {
		t.Errorf("incorrect first values was set. %#v", tc.FooFirst)
	}
	if tc.FooSecond != "FOO_SECOND" {
		t.Errorf("incorrect value was set. %#v", tc.FooSecond)
	}
}

func TestSimpleExternaEnvConfigTwoVariablesWithCommentOK(t *testing.T) {
	env := bytes.NewBuffer([]byte("# COMMENT=COMMENT\nFOO_FIRST=FOO_FIRST\nFOO_SECOND=FOO_SECOND\n"))
	tc := struct {
		Comment   string `env:"COMMENT"`
		FooFirst  string `env:"FOO_FIRST"`
		FooSecond string `env:"FOO_SECOND"`
	}{}
	envf := NewEnvConf()
	if err := envf.Parse(env); err != nil {
		t.Error("invalid parse env data:", err)
	}
	if err := envf.Set(); err != nil {
		t.Error("invalid set env variables:", err)
	}
	err := Parse(&tc)
	if err != nil {
		t.Error(err)
	}
	if tc.FooFirst != "FOO_FIRST" {
		t.Errorf("incorrect first values was set. %#v", tc.FooFirst)
	}
	if tc.FooSecond != "FOO_SECOND" {
		t.Errorf("incorrect value was set. %#v", tc.FooSecond)
	}
	if tc.Comment != "" {
		t.Errorf("incorrect values was set. %#v", tc.Comment)
	}
}

func TestSimpleExternaEnvConfigTwoVariablesWithQueteOK(t *testing.T) {
	env := bytes.NewBuffer([]byte(`FOO_FIRST="FOO_FIRST"
	FOO_SECOND='FOO_SECOND'`))
	tc := struct {
		FooFirst  string `env:"FOO_FIRST"`
		FooSecond string `env:"FOO_SECOND"`
	}{}
	envf := NewEnvConf()
	if err := envf.Parse(env); err != nil {
		t.Error("invalid parse env data:", err)
	}
	if err := envf.Set(); err != nil {
		t.Error("invalid set env variables:", err)
	}
	err := Parse(&tc)
	if err != nil {
		t.Error(err)
	}
	if tc.FooFirst != "FOO_FIRST" {
		t.Errorf("incorrect first values was set. %#v", tc.FooFirst)
	}
	if tc.FooSecond != "FOO_SECOND" {
		t.Errorf("incorrect value was set. %#v", tc.FooSecond)
	}
}

func TestSimpleExternaEnvConfigTwoVariablesEmptyStringOK(t *testing.T) {
	env := bytes.NewBuffer([]byte(`FOO_FIRST=""
	FOO_SECOND=''`))
	tc := struct {
		FooFirst  string `env:"FOO_FIRST"`
		FooSecond string `env:"FOO_SECOND"`
	}{
		FooFirst:  "foo_first",
		FooSecond: "foo_second",
	}
	envf := NewEnvConf()
	if err := envf.Parse(env); err != nil {
		t.Error("invalid parse env data:", err)
	}
	if err := envf.Set(); err != nil {
		t.Error("invalid set env variables:", err)
	}
	err := Parse(&tc)
	if err != nil {
		t.Error(err)
	}
	if tc.FooFirst != "" {
		t.Errorf("incorrect first values was set. %#v", tc.FooFirst)
	}
	if tc.FooSecond != "" {
		t.Errorf("incorrect value was set. %#v", tc.FooSecond)
	}
}

func TestParseEnvVars(t *testing.T) {
	tc := []string{
		"foo=bar",
		" foo   =bar ",
		"   foo    =   \"   bar   \"    ",
		"   foo    =   `bar`    ",
	}
	res := struct {
		Foo string `env:"foo"`
	}{}
	for _, c := range tc {
		envf := NewEnvConf()
		if err := envf.Parse(bytes.NewBufferString(c)); err != nil {
			t.Error("invalid parse env data:", err)
		}
		if err := envf.Set(); err != nil {
			t.Error("invalid set env variables:", err)
		}
		err := Parse(&res)
		if err != nil {
			t.Error(err)
		}
		if res.Foo != "bar" {
			t.Errorf("incorrect value was set. %#v", res.Foo)
		}
	}
}
