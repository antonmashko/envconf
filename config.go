package envconf

type Config interface {
	RawMessage() []byte
	Contains(keyName string) bool
	Unmarshal(data interface{}) error
}
