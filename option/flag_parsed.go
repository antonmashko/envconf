package option

type flagParsedFunc func() error

func (fP flagParsedFunc) Apply(opts *Options) {
	opts.flagParsed = fP
}

// WithFlagParsed define this callback when you need handle flags
// This callback will raise after method flag.Parse()
// return not nil error interrupt pasring
func WithFlagParsed(flagParsed func() error) ClientOption {
	return flagParsedFunc(flagParsed)
}
