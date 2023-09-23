package option

import (
	"errors"
	"testing"
)

func TestWithFlagParsed_Ok(t *testing.T) {
	cErr := errors.New("custom err")
	f := func() error {
		return cErr
	}
	opt := WithFlagParsed(f)
	opts := &Options{}
	opt.Apply(opts)
	err := opts.FlagParsed()()
	if err != cErr {
		t.Fatal("unexpected result: ", err)
	}
}
