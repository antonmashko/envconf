package envconf

import (
	"encoding/json"
)

type JsonConfig struct {
	data []byte
}

func NewJsonConfig() *JsonConfig {
	return &JsonConfig{}
}

func (j *JsonConfig) TagName() string {
	return "json"
}

func (j *JsonConfig) Read(data []byte) {
	j.data = data
}

func (j *JsonConfig) Unmarshal(v interface{}) error {
	if j.data == nil {
		return nil
	}
	return json.Unmarshal(j.data, v)
}
