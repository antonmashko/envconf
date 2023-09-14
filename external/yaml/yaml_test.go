package yaml

import (
	"reflect"
	"testing"
)

func TestYamlConf_ParseSimple_Ok(t *testing.T) {
	var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`
	tc := struct {
		A string
		B struct {
			RenamedC int   `yaml:"c"`
			D        []int `yaml:",flow"`
		}
	}{}
	extConf := NewYamlConf()
	extConf.Read([]byte(data))
	err := extConf.Unmarshal(&tc)
	if err != nil {
		t.Error("unexpected error")
	}

	if tc.A != "Easy!" && tc.B.RenamedC != 2 && !reflect.DeepEqual([]int{3, 4}, tc.B.D) {
		t.Errorf("incorrect values: %#v", tc)
	}
}
