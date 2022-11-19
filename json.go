package envconf

import (
	"encoding/json"
	"strings"
)

type JsonConfig struct {
	m    map[string]interface{}
	data []byte
}

func NewJsonConfig() *JsonConfig {
	return &JsonConfig{
		m: make(map[string]interface{}),
	}
}

func (j *JsonConfig) Read(data []byte) {
	j.data = data
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
			// lookup with ignore case
			for k, v := range mp {
				if strings.ToLower(k) == name {
					tmp = v
				}
			}
			if tmp == nil {
				// NOTE: not found
				return nil, false
			}
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
	if j.data == nil {
		return nil
	}
	err := json.Unmarshal(j.data, &j.m)
	if err != nil {
		return err
	}
	return json.Unmarshal(j.data, v)
}
