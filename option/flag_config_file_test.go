package option

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/antonmashko/envconf/external"
	"github.com/antonmashko/envconf/external/json"
)

func TestWithFlagConfigFile_Ok(t *testing.T) {
	f, err := os.CreateTemp("", "envconf.tmp")
	if err != nil {
		t.Fatal("os.Create:", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())
	const content = `{"foo":"bar"}`
	_, err = io.WriteString(f, content)
	if err != nil {
		t.Fatal("io.WriteString: ", err)
	}
	opt := WithFlagConfigFile("config1", f.Name(), "", func(b []byte) (external.External, error) {
		return json.Json(b), nil
	})
	opts := &Options{}
	opt.Apply(opts)
	if err = opts.FlagParsed()(); err != nil {
		t.Fatal("opts.FlagParsed(): ", err)
	}
	result := struct {
		Foo string `json:"foo"`
	}{}
	opts.External().Unmarshal(&result)
	if result.Foo != "bar" {
		t.Fatal("unexpected result: ", result)
	}
}

func TestWithFlagConfigFile_NotExist_Err(t *testing.T) {
	opt := WithFlagConfigFile("config2", "./conf.json", "", func(b []byte) (external.External, error) {
		return nil, nil
	})
	opts := &Options{}
	opt.Apply(opts)
	if err := opts.FlagParsed()(); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestWithFlagConfigFile_CustomError_Err(t *testing.T) {
	f, err := os.CreateTemp("", "envconf.tmp")
	if err != nil {
		t.Fatal("os.Create:", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	cErr := errors.New("custom error")
	opt := WithFlagConfigFile("config3", f.Name(), "", func(b []byte) (external.External, error) {
		return nil, cErr
	})
	opts := &Options{}
	opt.Apply(opts)
	if err = opts.FlagParsed()(); err == nil || err != cErr {
		t.Fatal("opts.FlagParsed() unexpected error: ", err)
	}
}

func TestWithFlagConfigFile_NilExternal_Ok(t *testing.T) {
	opt := &withExternalConfigFileOption{}
	if len(opt.TagName()) > 0 {
		t.Fatal("len(TagName()) greater than zero")
	}
	if opt.Unmarshal(nil) != nil {
		t.Fatal("opt.Unmarshal error not nil")
	}
}

func TestWithFlagConfigFile_JsonExternal_Ok(t *testing.T) {
	opt := &withExternalConfigFileOption{External: json.Json("{}")}
	if len(opt.TagName()) == 0 {
		t.Fatal("len(TagName()) is zero")
	}
	if opt.Unmarshal(&struct{}{}) != nil {
		t.Fatal("opt.Unmarshal error not nil")
	}
}
