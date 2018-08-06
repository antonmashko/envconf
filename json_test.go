package envconf

import "testing"

func TestSimpleExternalJsonConfigOK(t *testing.T) {
	json := `{"foo":"bar"}`
	tc := struct {
		Foo string `default:"fail"`
	}{}
	jconf := make(JsonConfig)
	if err := jconf.Unmarshal([]byte(json)); err != nil {
		t.Errorf("failed to unmarshal. err=%s", err)
	}
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
	jconf := make(JsonConfig)
	if err := jconf.Unmarshal([]byte(json)); err != nil {
		t.Errorf("failed to unmarshal. err=%s", err)
	}
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
	jconf := make(JsonConfig)
	if err := jconf.Unmarshal([]byte(json)); err != nil {
		t.Errorf("failed to unmarshal. err=%s", err)
	}
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
	jconf := make(JsonConfig)
	if err := jconf.Unmarshal([]byte(json)); err != nil {
		t.Errorf("failed to unmarshal. err=%s", err)
	}
	if err := ParseWithExternal(&tc, jconf); err != nil {
		t.Errorf("failed to external parse. err=%s", err)
	}
	if tc.FooBar.FooBar != "foo_bar" {
		t.Errorf("incorrect value was set. %#v", tc.FooBar)
	}
}
