package envconf

import (
	"os"
	"testing"
	"time"
)

func TestParseEmptyStructOK(t *testing.T) {
	if err := Parse(&struct{}{}); err != nil {
		t.Errorf("failed to parse empty struct. err=%#v", err)
	}
}

func TestNilDataNegative(t *testing.T) {
	if err := Parse(nil); err == nil || err != ErrNilData {
		t.Errorf("err doesn't equals to ErrNilData. err=%#v", err)
	}
}

func TestDoublePointerOK(t *testing.T) {
	tc := &struct {
		X string `default:"ok"`
	}{}
	if err := Parse(&tc); err != nil {
		t.Errorf("failed to parse. err=%s", err)
	}
	if tc.X != "ok" {
		t.Errorf("incorrect value was set. %#v", tc.X)
	}
}

func TestParsingDurtaionOK(t *testing.T) {
	tc := struct {
		X time.Duration `default:"5m"`
	}{}
	if err := Parse(&tc); err != nil {
		t.Errorf("failed to parse. err=%s", err)
	}
	if tc.X.Minutes() != 5.0 {
		t.Errorf("incorrect value was set. %#v", tc.X)
	}
}

func TestNilValueOK(t *testing.T) {
	tc := struct {
		X *string `default:"fail"`
	}{}
	IgnoreNilData = true
	if err := Parse(&tc); err != nil {
		t.Errorf("failed to parse. err=%s", err)
	}
	if tc.X != nil {
		t.Errorf("incorrect value was set. %#v", *tc.X)
	}
}

func TestFlagParsedCallbackOK(t *testing.T) {
	x := 0
	FlagParsed = func() error {
		x = 1
		return nil
	}
	tc := struct{}{}
	if err := Parse(&tc); err != nil {
		t.Errorf("failed to parse. err=%s", err)
	}
	if x != 1 {
		t.Errorf("incorrect value was set. %#v", x)
	}
}

type foo struct {
	Bar bar
}

type bar struct {
	Data string `env:"data-env"`
}

func TestHotReload(t *testing.T) {
	ts := &foo{}
	f := func(b *bar) {
		t.Logf("%p %s", b, b.Data)
	}
	os.Setenv("data-env", "1")
	Parse(&ts)
	f(&ts.Bar)
	os.Setenv("data-env", "2")
	Parse(&ts)
	f(&ts.Bar)
}
