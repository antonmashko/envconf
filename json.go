package envconf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"strings"
)

const JsonConfigName = "json"

type JsonConfig struct {
	filepath string
	values   map[string]interface{}
}

func (j *JsonConfig) SetFilePathFlag(flagName string, defaultPath string) {
	flag.StringVar(&j.filepath, flagName, defaultPath, "json config file path")
}

func (j *JsonConfig) RawMessage() []byte {
	if j.filepath == "" {
		elog.Infoln("envconf: json filepath is not defined")
		return nil
	}
	b, err := ioutil.ReadFile(j.filepath)
	if err != nil {
		elog.Errorf("envconf: failed to read file=%s error=%s", j.filepath, err)
		return nil
	}
	return b
}

func (j *JsonConfig) Contains(keyName string) bool {
	values := strings.Split(keyName, string(Separator))
	m := j.values
	for i, v := range values {
		tm, ok := m[strings.ToLower(v)]
		if !ok {
			return false
		}
		if i+1 == len(values) {
			break
		}
		m, ok = tm.(map[string]interface{})
		if !ok {
			return false
		}
	}
	return true
}

func (j *JsonConfig) Unmarshal(data interface{}) error {
	b := j.RawMessage()
	if b == nil {
		return nil
	}
	if err := json.Unmarshal(b, data); err != nil {
		return err
	}
	return json.Unmarshal(b, &j.values)
}
