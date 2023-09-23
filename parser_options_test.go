package envconf_test

import (
	"os"
	"testing"

	"github.com/antonmashko/envconf"
	"github.com/antonmashko/envconf/option"
)

func TestPriority_FromFlag_Ok(t *testing.T) {
	const expected = "from-flag-value"
	os.Args = append(os.Args, "-ptest-field1="+expected)
	os.Clearenv()
	os.Setenv("TEST_FIELD", "from-env-value")
	data := struct {
		Field string `flag:"ptest-field1" env:"TEST_FIELD" default:"default-variable"`
	}{}
	ecfg := envconf.New()
	if err := ecfg.Parse(&data, option.WithPriorityOrder(option.FlagVariable, option.EnvVariable, option.DefaultValue)); err != nil {
		t.Fatal("Parse: ", err)
	}
	if data.Field != expected {
		t.Fatalf("incorrect result. expected=%s actual=%s", expected, data.Field)
	}
}

func TestPriority_FromEnv_Ok(t *testing.T) {
	const expected = "from-env-value"
	os.Args = append(os.Args, "-ptest-field2=from-flag-value")
	os.Clearenv()
	os.Setenv("TEST_FIELD", expected)
	data := struct {
		Field string `flag:"ptest-field2" env:"TEST_FIELD" default:"default-variable"`
	}{}
	ecfg := envconf.New()
	if err := ecfg.Parse(&data, option.WithPriorityOrder(option.EnvVariable, option.FlagVariable, option.DefaultValue)); err != nil {
		t.Fatal("Parse: ", err)
	}
	if data.Field != expected {
		t.Fatalf("incorrect result. expected=%s actual=%s", expected, data.Field)
	}
}

func TestPriority_FromDefault_Ok(t *testing.T) {
	const expected = "default-variable"
	os.Args = append(os.Args, "-ptest-field3=from-flag-value")
	os.Clearenv()
	os.Setenv("TEST_FIELD", "from-env-value")
	data := struct {
		Field string `flag:"ptest-field3" env:"TEST_FIELD" default:"default-variable"`
	}{}
	ecfg := envconf.New()
	if err := ecfg.Parse(&data, option.WithPriorityOrder(option.DefaultValue, option.FlagVariable, option.EnvVariable)); err != nil {
		t.Fatal("Parse: ", err)
	}

	if data.Field != expected {
		t.Fatalf("incorrect result. expected=%s actual=%s", expected, data.Field)
	}
}

func TestPriority_EmptyPriorityDefineFromFlag_Ok(t *testing.T) {
	const expected = "from-flag-value"
	os.Args = append(os.Args, "-ptest-field4=from-flag-value")
	os.Clearenv()
	os.Setenv("TEST_FIELD", "from-env-value")
	data := struct {
		Field string `flag:"ptest-field4" env:"TEST_FIELD" default:"default-variable"`
	}{}
	ecfg := envconf.New()
	if err := ecfg.Parse(&data, option.WithPriorityOrder()); err != nil {
		t.Fatal("Parse: ", err)
	}

	if data.Field != expected {
		t.Fatalf("incorrect result. expected=%s actual=%s", expected, data.Field)
	}
}

func TestPriority_InvalidConfigSourceDefineFromFlag_Ok(t *testing.T) {
	const expected = "from-flag-value"
	os.Args = append(os.Args, "-ptest-field5=from-flag-value")
	os.Clearenv()
	os.Setenv("TEST_FIELD", "from-env-value")
	data := struct {
		Field string `flag:"ptest-field5" env:"TEST_FIELD" default:"default-variable"`
	}{}
	ecfg := envconf.New()
	if err := ecfg.Parse(&data, option.WithPriorityOrder(option.ConfigSource(123), option.ConfigSource(124))); err != nil {
		t.Fatal("Parse: ", err)
	}
	if data.Field != expected {
		t.Fatalf("incorrect result. expected=%s actual=%s", expected, data.Field)
	}
}
