package yaml

import "gopkg.in/yaml.v3"

type YamlConf struct {
	data []byte
}

func NewYamlConf() *YamlConf {
	return &YamlConf{}
}

func (c *YamlConf) Read(data []byte) {
	c.data = data
}

func (c *YamlConf) TagName() string {
	return "yaml"
}

func (c *YamlConf) Unmarshal(v interface{}) error {
	return yaml.Unmarshal(c.data, v)
}
