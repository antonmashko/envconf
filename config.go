package envconf

var configs map[string]Config

type Config interface {
	Contains(keyName string) bool
	Unmarshal(filepath string, data interface{}) error
}

func Registre(name string, config Config) {
	configs[name] = config
}
