package envconf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"reflect"
	"strings"
)

const JsonConfigName = "json"

type JsonConfig struct {
	filepath string
	values   JsonConfigValues
}

type JsonConfigValues map[string]interface{}

func (j *JsonConfig) SetFilePathFlag(flagName string, defaultPath string) {
	flag.StringVar(&j.filepath, flagName, defaultPath, "json config file path")
}

func (j *JsonConfig) RawMessage() []byte {
	if j.filepath == "" {
		traceLogger.Println("envconf: json filepath is not defined")
		return nil
	}
	b, err := ioutil.ReadFile(j.filepath)
	if err != nil {
		errorLogger.Printf("envconf: failed to read file=%s error=%s", j.filepath, err)
		return nil
	}
	return b
}

func (j *JsonConfig) Contains(keyName string) bool {
	values := strings.Split(keyName, string(Separator))
	m := j.values
	for i, v := range values {
		tm, ok := m[v]
		if !ok {
			return false
		}
		if i+1 == len(values) {
			break
		}
		m, ok = tm.(JsonConfigValues)
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

	j.values = j.createJsonConfigValues(reflect.ValueOf(data))

	return nil
}

func (j *JsonConfig) createJsonConfigValues(val reflect.Value) (values JsonConfigValues) {
	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	values = make(JsonConfigValues)

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
			elm := valueField.Elem()
			if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
				valueField = elm
			}
		}

		if valueField.Kind() == reflect.Ptr {
			valueField = valueField.Elem()

		}

		if valueField.Kind() == reflect.Struct {
			values[typeField.Name] = j.createJsonConfigValues(valueField)
		} else {
			values[typeField.Name] = valueField
		}
	}
	return
}
