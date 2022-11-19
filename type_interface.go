package envconf

type interfaceType struct{}

func (t *interfaceType) Init() error {
	// not supported
	return nil
}

func (t *interfaceType) Define() error {
	return nil
}
