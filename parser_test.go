package envconf

import (
	"log"
	"os"
	"strconv"
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

func TestParseStringEnvVariableOK(t *testing.T) {
	const v = "TEST"
	os.Setenv(v, v)
	tc := &struct {
		X string `env:"TEST"`
	}{}
	err := Parse(&tc)
	if err != nil {
		t.Error("failed to parse. err=", err)
	}
	if tc.X != v {
		t.Errorf("incorrect result of env variable='%s'", tc.X)
	}
}

func TestParseIntEnvVariableOK(t *testing.T) {
	const v = "TEST"
	const result = 10
	os.Setenv(v, strconv.Itoa(result))
	tc := &struct {
		X int `env:"TEST"`
	}{}
	err := Parse(&tc)
	if err != nil {
		t.Error("failed to parse. err=", err)
	}
	if tc.X != result {
		t.Errorf("incorrect result of env variable='%d'", tc.X)
	}
}

func TestRequiredFieldOK(t *testing.T) {
	tc := &struct {
		X string `required:"true"`
	}{}
	if err := Parse(&tc); err == nil {
		t.Error("expected error was not throw.")
	} else {
		t.Logf("%#[1]v %[1]T", err)
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
		t.Errorf("incorrect value was set. %#v", tc.X)
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

func TestSetLoggerOK(t *testing.T) {
	logger := &log.Logger{}
	SetLogger(logger)
	if debugLogger != logger {
		t.Errorf("Can't set logger")
	}
}
