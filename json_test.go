package envconf_test

import (
	"testing"

	"github.com/antonmashko/envconf"
)

func TestJsonConfig_SimpleExternalJsonConfig_OK(t *testing.T) {
	json := `{"foo":"bar"}`
	tc := struct {
		Foo string `default:"fail"`
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.Foo != "bar" {
		t.Errorf("incorrect value was set. %#v", tc.Foo)
	}
}

func TestJsonConfig_SimpleExternalFieldWithUnderscore_OK(t *testing.T) {
	json := `{"foo_bar":"foo_bar"}`
	tc := struct {
		FooBar string `json:"foo_bar" default:"fail"`
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.FooBar != "foo_bar" {
		t.Errorf("incorrect value was set. %#v", tc.FooBar)
	}
}

func TestJsonConfig_NestedStructExternal_OK(t *testing.T) {
	json := `{
		"foo": {
			"bar": {
				"foobar": "foo_bar"
			}
		}
	}`
	tc := struct {
		Foo struct {
			Bar struct {
				FooBar string `default:"fail"`
			}
		}
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.Foo.Bar.FooBar != "foo_bar" {
		t.Errorf("incorrect value was set. %#v", tc.Foo)
	}
}

func TestJsonConfig_NestedStructExternalFieldWithUnderscore_OK(t *testing.T) {
	json := `{
		"foo_bar": {
			"foo_bar": "foo_bar"
		}
	}`
	tc := struct {
		FooBar struct {
			FooBar string `json:"foo_bar" default:"fail"`
		} `json:"foo_bar"`
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.FooBar.FooBar != "foo_bar" {
		t.Errorf("incorrect value was set. %#v", tc.FooBar)
	}
}

func TestJsonConfig_Slice_OK(t *testing.T) {
	json := `{
		"foo": [
			1
		]
	}`
	tc := struct {
		Foo []int
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if len(tc.Foo) != 1 || tc.Foo[0] != 1 {
		t.Errorf("incorrect value was set. %#v", tc.Foo)
	}
}

func TestJsonConfig_SliceFloat_Ok(t *testing.T) {
	json := `{
		"foo": [
			1.1
		]
	}`
	tc := struct {
		Foo []float32
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if len(tc.Foo) != 1 || tc.Foo[0] != 1.1 {
		t.Errorf("incorrect value was set. %#v", tc.Foo)
	}
}

func TestJsonConfig_PropertyCamelCase_Ok(t *testing.T) {
	json := `{
		"Foo": {
			"Bar": {
				"FooBar": "foo_bar"
			}
		}
	}`
	tc := struct {
		Foo struct {
			Bar struct {
				FooBar string `default:"fail"`
			}
		}
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.Foo.Bar.FooBar != "foo_bar" {
		t.Errorf("incorrect value was set. %#v", tc.Foo)
	}
}

func TestJsonConfig_CaseSensitive_Ok(t *testing.T) {
	json := `{
		"abc": 1,
		"Abc": 2,
		"ABC": {
			"abc": 3
		}
	}`
	tc := struct {
		AbC int `json:"abc"`
		Abc int
		ABC struct {
			ABC int
		}
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.AbC != 1 || tc.Abc != 2 || tc.ABC.ABC != 3 {
		t.Errorf("incorrect value was set. %#v", tc)
	}
}

func TestJsonConfig_NonExistJsonValueDefaultUse_Ok(t *testing.T) {
	json := `{"foo":2}`
	tc := struct {
		Foo int `json:"foo"`
		Bar int `json:"bar" default:"5"`
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.Foo != 2 || tc.Bar != 5 {
		t.Errorf("incorrect value was set. %#v", tc)
	}
}

func TestJsonConfig_NonExistConfigValue_Ok(t *testing.T) {
	json := `{"foo":2, "bar":5}`
	tc := struct {
		Foo int `json:"foo"`
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.Foo != 2 {
		t.Errorf("incorrect value was set. %#v", tc)
	}
}

func TestJsonConfig_IncorrectType_Err(t *testing.T) {
	json := `{"foo":2, "bar":{"abc":3}}`
	tc := struct {
		Foo int `json:"foo"`
		Bar int `json:"bar"`
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestJsonConfig_Map_Err(t *testing.T) {
	json := `{"foo":{"a":"b", "b":"c"}}`
	tc := struct {
		Foo map[string]string `json:"foo"`
	}{}
	jconf := envconf.NewJsonConfig()
	jconf.Read([]byte(json))
	if err := envconf.ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}

	if tc.Foo["a"] != "b" && tc.Foo["b"] != "c" {
		t.Errorf("incorrect result: %#v", tc.Foo)
	}
}
