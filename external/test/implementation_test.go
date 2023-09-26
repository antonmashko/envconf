package test

import (
	"github.com/antonmashko/envconf"
	"github.com/antonmashko/envconf/external/yaml"
)

var _ envconf.External = (yaml.Yaml)([]byte{})
