package external

import (
	"testing"

	"github.com/antonmashko/envconf/option"
)

func TestWithFlagConfigFile_Ok(t *testing.T) {
	opt := WithFlagConfigFile("config", "", "", func(b []byte) error {
		return nil
	})
	opts := &option.Options{}
	opt.Apply(opts)
}
