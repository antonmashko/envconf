package envconf

import (
	"encoding/json"
)

// Json implementation of External Configuration source
type Json []byte

func (j Json) TagName() string {
	return "json"
}

func (j Json) Unmarshal(v interface{}) error {
	return json.Unmarshal(j, v)
}
