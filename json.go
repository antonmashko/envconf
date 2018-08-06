package envconf

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JsonConfig map[string]interface{}

func (j JsonConfig) Get(values ...Value) (string, bool) {
	const tagName = "json"
	mp := map[string]interface{}(j)
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
			return fmt.Sprint(tmp), true
		}
	}
	return "", false
}

func (j JsonConfig) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &j)
}
