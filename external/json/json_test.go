package json

import (
	"reflect"
	"testing"
)

func TestJson_ParseSimple_Ok(t *testing.T) {
	const json = `
{
	"a":1,
	"b": [1, 2],
	"c": {
		"a":"3"
	}
}
	`
	jsonConf := Json([]byte(json))
	if !reflect.DeepEqual(jsonConf.TagName(), []string{"json"}) {
		t.Fatal("tag:", jsonConf.TagName())
	}
	tc := struct {
		A int
		B []uint
		C struct {
			A string
		}
	}{}
	err := jsonConf.Unmarshal(&tc)
	if err != nil {
		t.Fatal("jsonConf.Unmarshal: ", err)
	}
	if tc.A != 1 || !reflect.DeepEqual(tc.B, []uint{1, 2}) || tc.C.A != "3" {
		t.Fatalf("incorrect result: %#v", tc)
	}
}
