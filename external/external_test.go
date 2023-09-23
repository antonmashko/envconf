package external

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/antonmashko/envconf/option"
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
	var result string
	opt := WithFlagConfigFile("config1", f.Name(), "", func(b []byte) error {
		result = string(b)
		return nil
	})
	opts := &option.Options{}
	opt.Apply(opts)
	if err = opts.FlagParsed()(); err != nil {
		t.Fatal("opts.FlagParsed(): ", err)
	}
	if result != content {
		t.Fatal("unexpected result: ", result)
	}
}

func TestWithFlagConfigFile_NotExist_Err(t *testing.T) {
	opt := WithFlagConfigFile("config2", "./conf.json", "", func(b []byte) error {
		return nil
	})
	opts := &option.Options{}
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
	opt := WithFlagConfigFile("config3", f.Name(), "", func(b []byte) error {
		return cErr
	})
	opts := &option.Options{}
	opt.Apply(opts)
	if err = opts.FlagParsed()(); err == nil || err != cErr {
		t.Fatal("opts.FlagParsed() unexpected error: ", err)
	}
}
