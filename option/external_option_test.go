package option

import (
	"testing"

	"github.com/antonmashko/envconf/external/json"
)

func TestExternalOption_Ok(t *testing.T) {
	jsonConf := json.Json("{}")
	opts := &Options{}
	WithExternal(jsonConf).Apply(opts)
	if opts.External() == nil {
		t.Fatal("unexpected result")
	}
}
