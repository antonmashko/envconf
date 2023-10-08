package external_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/antonmashko/envconf"
	jsonconf "github.com/antonmashko/envconf/external/json"
	"github.com/antonmashko/envconf/option"
)

func TestJsonConfig_SimpleExternalJsonConfig_OK(t *testing.T) {
	json := `{"foo":"bar"}`
	tc := struct {
		Foo string `default:"fail"`
	}{}
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	const expected = "foo_bar"
	if tc.Foo.Bar.FooBar != expected {
		t.Errorf("incorrect value. expected=%s actual=%s", expected, tc.Foo.Bar.FooBar)
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
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
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err == nil {
		t.Errorf("expected error but got nil")
	}
}

func TestJsonConfig_Array_Ok(t *testing.T) {
	json := `{"foo":[2, 3, 4, 5]}`
	tc := struct {
		Foo [4]int `json:"foo"`
	}{}
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}

	if tc.Foo[0] != 2 && tc.Foo[1] != 3 && tc.Foo[2] != 4 && tc.Foo[3] != 5 {
		t.Errorf("incorrect result: %#v", tc)
	}
}

func TestJsonConfig_Map_Ok(t *testing.T) {
	json := `{"foo":{"a":"b", "b":"c"}}`
	tc := struct {
		Foo map[string]string `json:"foo"`
	}{}
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}

	if tc.Foo["a"] != "b" && tc.Foo["b"] != "c" {
		t.Errorf("incorrect result: %#v", tc)
	}
}

func TestJsonConfig_TagWithExtra_Ok(t *testing.T) {
	json := `{"bar":123}`
	tc := struct {
		Foo int `json:"bar,omitempty"`
	}{}
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}

	if tc.Foo != 123 {
		t.Errorf("incorrect result: %#v", tc)
	}
}

func TestJsonConfig_ZeroValue_Ok(t *testing.T) {
	json := `{"foo":0}`
	tc := struct {
		Foo int `default:"5"`
	}{}
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.Foo != 0 {
		t.Errorf("incorrect result: %#v", tc)
	}
}

func TestJsonConfig_ValueWithSpace_Ok(t *testing.T) {
	json := `{"foo": "   bar "}`
	tc := struct {
		Foo string
	}{}
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}

	if tc.Foo != "   bar " {
		t.Errorf("incorrect result: %#v", tc)
	}
}

func TestJsonConfig_SliceOfStructs_Ok(t *testing.T) {
	json := `{"foo": [
		{"f1":"bar1"}
	]}`
	tc := struct {
		Foo []struct {
			Field1 string `json:"f1"`
		}
	}{}
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}

	if len(tc.Foo) < 1 || tc.Foo[0].Field1 != "bar1" {
		t.Errorf("incorrect result: %#v", tc)
	}
}

func TestJsonConfig_SliceOfStructsComplex_Ok(t *testing.T) {
	json := `{
		"foo": [{
			"f1": "bar1",
			"sl1": [
				"test",
				1,
				42.2,
				true
			]
		}]
	}`
	tc := struct {
		Foo []struct {
			Field1 string        `json:"f1"`
			Sl     []interface{} `json:"sl1"`
		}
	}{}
	if err := envconf.Parse(&tc, option.WithExternal(jsonconf.Json([]byte(json)))); err != nil {
		t.Fatalf("failed to external parse. err=%s", err)
	}

	if len(tc.Foo) < 1 || tc.Foo[0].Field1 != "bar1" {
		t.Fatalf("incorrect result: %#v", tc)
	}

	expectedSl := []interface{}{"test", float64(1), 42.2, true}
	if !reflect.DeepEqual(expectedSl, tc.Foo[0].Sl) {
		t.Fatalf("incorrect result. expected=%v actual=%v", expectedSl, tc.Foo[0].Sl)
	}
}

func TestJsonConfig_EnvVarInjection_Ok(t *testing.T) {
	json := `{
		"foo": {
			"bar": "${.env.JSON_TEST_ENVVAR_INJECTION}"
		}
	}`
	tc := struct {
		Foo struct {
			Bar string `json:"bar"`
		}
	}{}

	expectedValue := "foo_bar"
	os.Setenv("JSON_TEST_ENVVAR_INJECTION", expectedValue)
	err := envconf.Parse(&tc,
		option.WithExternal(jsonconf.Json([]byte(json))),
		option.WithExternalEnvInjection(true),
	)
	if err != nil {
		t.Fatalf("failed to external parse. err=%s", err)
	}

	if tc.Foo.Bar != expectedValue {
		t.Fatalf("incorrect result. expected=%s actual=%s", expectedValue, tc.Foo.Bar)
	}
}

func TestJsonConfig_EnvVarInjectionFromCollection_Ok(t *testing.T) {
	json := `{
		"foo": [{
			"bar": "${.env.JSON_TEST_ENVVAR_INJECTION}"
		}]
	}`
	tc := struct {
		Foo []struct {
			Bar string `json:"bar"`
		}
	}{}

	expectedValue := "foo_bar"
	os.Setenv("JSON_TEST_ENVVAR_INJECTION", expectedValue)
	err := envconf.Parse(&tc,
		option.WithExternal(jsonconf.Json([]byte(json))),
		option.WithExternalEnvInjection(true),
	)
	if err != nil {
		t.Fatalf("failed to external parse. err=%s", err)
	}

	if tc.Foo[0].Bar != expectedValue {
		t.Fatalf("incorrect result. expected=%s actual=%s", expectedValue, tc.Foo[0].Bar)
	}
}
