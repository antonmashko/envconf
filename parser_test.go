package envconf_test

import (
	"testing"

	"github.com/antonmashko/envconf"
)

func TestParse_EmptyStruct_OK(t *testing.T) {
	if err := envconf.Parse(&struct{}{}); err != nil {
		t.Errorf("failed to parse empty struct. err=%#v", err)
	}
}

func TestParse_NilData_NilDataErr(t *testing.T) {
	if err := envconf.Parse(nil); err == nil || err != envconf.ErrNilData {
		t.Errorf("err doesn't equals to ErrNilData. err=%#v", err)
	}
}

func TestParse_NilValue_NilDataErr(t *testing.T) {
	var d *struct{}
	if err := envconf.Parse(d); err == nil || err != envconf.ErrNilData {
		t.Errorf("err doesn't equals to ErrNilData. err=%#v", err)
	}
}

func TestParse_PassDataByValue_Err(t *testing.T) {
	data := struct {
		Field string `default:"123"`
	}{}
	if err := envconf.Parse(data); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestParse_FlagParsedCallback_OK(t *testing.T) {
	x := 0
	envconf.FlagParsed = func() error {
		x = 1
		return nil
	}
	tc := struct{}{}
	if err := envconf.Parse(&tc); err != nil {
		t.Errorf("failed to parse. err=%s", err)
	}
	if x != 1 {
		t.Errorf("incorrect value was set. %#v", x)
	}
}

func TestParse_InvalidData_Err(t *testing.T) {
	var cfg string
	if err := envconf.Parse(&cfg); err == nil {
		t.Fail()
	}
}
