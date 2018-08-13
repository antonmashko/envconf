package envconf

import (
	"encoding/json"
	"strings"
)

type JsonConfig struct {
	m    map[string]interface{}
	data []byte
}

func NewJsonConfig(jsonData []byte) *JsonConfig {
	return &JsonConfig{
		m:    make(map[string]interface{}),
		data: jsonData,
	}
}

func (j *JsonConfig) Get(values ...Value) (interface{}, bool) {
	const tagName = "json"
	mp := map[string]interface{}(j.m)
	for _, v := range values {
		name := v.Tag().Tag.Get(tagName)
		if name == "" {
			name = v.Name()
		}
		name = strings.ToLower(name)
		tmp, ok := mp[name]
		if !ok {
			return "", false
		}
		switch tmp.(type) {
		case map[string]interface{}:
			mp = tmp.(map[string]interface{})
			break
		default:
			return tmp, true
		}
	}
	return nil, false
}

func (j *JsonConfig) Unmarshal(v interface{}) error {
	return json.Unmarshal(j.data, &j.m)
}
