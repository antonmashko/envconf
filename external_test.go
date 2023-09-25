package envconf

import (
	"testing"
)

func TestExternal_EmpExt_Ok(t *testing.T) {
	e := emptyExt{}
	if e.TagName() != "" || e.Unmarshal(nil) != nil {
		t.Error("unexpected result")
	}
}

func TestExternal_newExternalConfig_Ok(t *testing.T) {
	ext := newExternalConfig(emptyExt{})
	if ext == nil {
		t.Fail()
	}
	if ext.Unmarshal(nil) != nil {
		t.Error("unexpected result")
	}
}

func TestExternal_InvalidJson_Err(t *testing.T) {
	jsonConf := Json([]byte("<test></test>"))
	ext := newExternalConfig(jsonConf)
	tc := struct {
		Foo int
	}{}
	if ext.Unmarshal(&tc) == nil {
		t.Error("unexpected error got nil")
	}
}
