package yaml

import "gopkg.in/yaml.v3"

type Yaml []byte

func (y Yaml) TagName() []string {
	return []string{"yaml"}
}

func (y Yaml) Unmarshal(v interface{}) error {
	return yaml.Unmarshal(y, v)
}
