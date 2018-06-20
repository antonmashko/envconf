package envconf

type Config interface {
	RawMessage() []byte
	Contains(keyName string) bool
	Unmarshal(data interface{}) error
}

type emptyConfig struct{}

func (c *emptyConfig) RawMessage() []byte { return nil }

func (c *emptyConfig) Contains(keyName string) bool { return false }

func (c *emptyConfig) Unmarshal(data interface{}) error { return nil }
