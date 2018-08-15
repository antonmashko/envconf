package envconf

import "testing"

func TestSimpleExternalJsonConfigOK(t *testing.T) {
	json := `{"foo":"bar"}`
	tc := struct {
		Foo string `default:"fail"`
	}{}
	jconf := NewJsonConfig()
	jconf.Read([]byte(json))
	if err := ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.Foo != "bar" {
		t.Errorf("incorrect value was set. %#v", tc.Foo)
	}
}

func TestSimpleExternalJsonConfigFieldWithUnderscoreOK(t *testing.T) {
	json := `{"foo_bar":"foo_bar"}`
	tc := struct {
		FooBar string `json:"foo_bar" default:"fail"`
	}{}
	jconf := NewJsonConfig()
	jconf.Read([]byte(json))
	if err := ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.FooBar != "foo_bar" {
		t.Errorf("incorrect value was set. %#v", tc.FooBar)
	}
}

func TestNestedStructExternalJsonConfigOK(t *testing.T) {
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
	jconf := NewJsonConfig()
	jconf.Read([]byte(json))
	if err := ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.Foo.Bar.FooBar != "foo_bar" {
		t.Errorf("incorrect value was set. %#v", tc.Foo)
	}
}

func TestNestedStructExternalJsonConfigFieldWithUnderscoreOK(t *testing.T) {
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
	jconf := NewJsonConfig()
	jconf.Read([]byte(json))
	if err := ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.FooBar.FooBar != "foo_bar" {
		t.Errorf("incorrect value was set. %#v", tc.FooBar)
	}
}

func TestSliceJsonConfigOK(t *testing.T) {
	json := `{
		"foo": [
			1
		]
	}`
	tc := struct {
		Foo []int
	}{}
	jconf := NewJsonConfig()
	jconf.Read([]byte(json))
	if err := ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if len(tc.Foo) != 1 || tc.Foo[0] != 1 {
		t.Errorf("incorrect value was set. %#v", tc.Foo)
	}
}

func TestSliceJsonConfigFloatOk(t *testing.T) {
	json := `{
		"foo": [
			1.1
		]
	}`
	tc := struct {
		Foo []float32
	}{}
	jconf := NewJsonConfig()
	jconf.Read([]byte(json))
	if err := ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if len(tc.Foo) != 1 || tc.Foo[0] != 1.1 {
		t.Errorf("incorrect value was set. %#v", tc.Foo)
	}
}

// func TestSliceJsonConfigNegative(t *testing.T) {
// 	json := `{
// 		"foo": [
// 			1.5
// 		]
// 	}`
// 	tc := struct {
// 		Foo []int
// 	}{}
// 	jconf := NewJsonConfig()
// 	jconf.Read([]byte(json))
// 	if err := ParseWithExternal(&tc, jconf); err != nil {
// 		t.Errorf("failed to external parse. err=%[1]s err_type=%[1]T", err)
// 	}
// 	if len(tc.Foo) != 1 || tc.Foo[0] != 0 {
// 		t.Errorf("incorrect value was set. %#v", tc.Foo)
// 	}
// }
