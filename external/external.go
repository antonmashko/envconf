package external

import (
	"flag"
	"fmt"
	"os"

	"github.com/antonmashko/envconf/option"
)

// WithFlagConfigFile wraps option.WithFlagParsed with reading configuration file from flag defined path
func WithFlagConfigFile(flagName string, flagValue string, flagDescription string, initConf func([]byte) error) option.ClientOption {
	cfg := flag.String(flagName, flagValue, flagDescription)
	return option.WithFlagParsed(func() error {
		b, err := os.ReadFile(*cfg)
		if err != nil {
			return fmt.Errorf("os.ReadFile: %w", err)
		}
		return initConf(b)
	})
}
