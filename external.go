package envconf

type External interface {
	Get(Value) (string, bool)
}

type emptyConfig struct{}

func (c *emptyConfig) Get(v Value) (string, bool) { return "", false }
