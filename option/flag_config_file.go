package option

import (
	"flag"
	"fmt"
	"os"

	"github.com/antonmashko/envconf/external"
)

type withExternalConfigFileOption struct {
	external.External
	fpOpt flagParsedFunc
}

func (o *withExternalConfigFileOption) TagName() []string {
	if o.External == nil {
		return []string{}
	}
	return o.External.TagName()
}

func (o *withExternalConfigFileOption) Unmarshal(v interface{}) error {
	if o.External == nil {
		return nil
	}
	return o.External.Unmarshal(v)
}

func (o *withExternalConfigFileOption) Apply(opts *Options) {
	o.fpOpt.Apply(opts)
	opts.external = o
}

// WithFlagConfigFile wraps option.WithFlagParsed with reading configuration file from flag defined path
func WithFlagConfigFile(flagName string, flagValue string, flagDescription string, initConf func([]byte) (external.External, error)) ClientOption {
	cfg := flag.String(flagName, flagValue, flagDescription)
	opt := &withExternalConfigFileOption{}
	opt.fpOpt = func() error {
		b, err := os.ReadFile(*cfg)
		if err != nil {
			return fmt.Errorf("os.ReadFile: %w", err)
		}
		ext, err := initConf(b)
		if err != nil {
			return err
		}
		opt.External = ext
		return nil
	}
	return opt
}
