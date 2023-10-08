package option

import (
	"testing"

	"github.com/antonmashko/envconf/external/json"
)

func TestConfigSource_Ok(t *testing.T) {
	cs := []ConfigSource{NoConfigValue, FlagVariable, EnvVariable, ExternalSource, DefaultValue}
	for i := range cs {
		t.Logf("%[1]s %[1]d", cs[i])
	}
}

func TestExternalOption_Ok(t *testing.T) {
	jsonConf := json.Json("{}")
	opts := &Options{}
	WithExternal(jsonConf).Apply(opts)
	if opts.External() == nil {
		t.Fatal("unexpected result")
	}
}
