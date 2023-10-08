package envconf_test

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/antonmashko/envconf"
)

func TestParse_Array_Ok(t *testing.T) {
	cfg := struct {
		Field [5]int `env:"TEST_PARSE_ARRAY_OK"`
	}{}
	os.Setenv("TEST_PARSE_ARRAY_OK", "-2, -1,0, 1 ,2 ")
	expectedResult := [5]int{-2, -1, 0, 1, 2}
	if err := envconf.Parse(&cfg); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Field, expectedResult) {
		t.Fatalf("incorrect result. expected=%v actual=%v", expectedResult, cfg.Field)
	}
}

func TestParse_Array_ErrOutOfRange(t *testing.T) {
	cfg := struct {
		Field [2]int `env:"TEST_PARSE_ARRAY_OK"`
	}{}
	os.Setenv("TEST_PARSE_ARRAY_OK", "-2, -1,0, 1 ,2 ")
	if err := envconf.Parse(&cfg); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestParse_Slice_Ok(t *testing.T) {
	cfg := struct {
		Field []int `env:"TEST_PARSE_SLICE_OK"`
	}{}
	os.Setenv("TEST_PARSE_SLICE_OK", "-2, -1,0, 1 ,2 ")
	expectedResult := []int{-2, -1, 0, 1, 2}
	if err := envconf.Parse(&cfg); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Field, expectedResult) {
		t.Fatalf("incorrect result. expected=%v actual=%v", expectedResult, cfg.Field)
	}
}

func TestParse_EmptyInterfaceSlice_Ok(t *testing.T) {
	cfg := struct {
		Slice []interface{}
	}{}
	if err := envconf.Parse(&cfg); err != nil {
		t.Fatal(err)
	}
}

func TestParse_InterfaceSlice_Ok(t *testing.T) {
	f1 := &struct {
		Foo string `env:"TEST_PARSE_INTERFACE_SLICE_OK_FOO1"`
	}{}
	f2 := &struct {
		Foo int `env:"TEST_PARSE_INTERFACE_SLICE_OK_FOO2"`
	}{}
	f3 := &struct {
		Foo float64 `env:"TEST_PARSE_INTERFACE_SLICE_OK_FOO3"`
	}{}
	cfg := struct {
		Field []interface{}
	}{
		Field: []interface{}{
			f1,
			f2,
			[]interface{}{f3},
		},
	}
	expectedF1 := "test"
	expectedF2 := 432
	expectedF3 := 43.2
	os.Setenv("TEST_PARSE_INTERFACE_SLICE_OK_FOO1", expectedF1)
	os.Setenv("TEST_PARSE_INTERFACE_SLICE_OK_FOO2", strconv.Itoa(expectedF2))
	os.Setenv("TEST_PARSE_INTERFACE_SLICE_OK_FOO3", fmt.Sprint(expectedF3))

	if err := envconf.Parse(&cfg); err != nil {
		t.Fatal(err)
	}
	if f1.Foo != expectedF1 {
		t.Fatalf("incorrect result. expected[0]=%v actual[0]=%v", expectedF1, f1.Foo)
	}
	if f2.Foo != expectedF2 {
		t.Fatalf("incorrect result. expected[1]=%v actual[1]=%v", expectedF2, f2.Foo)
	}
	if f3.Foo != expectedF3 {
		t.Fatalf("incorrect result. expected[1]=%v actual[1]=%v", expectedF3, f3.Foo)
	}
}

func TestParse_Slice_ErrInvalidElement(t *testing.T) {
	cfg := struct {
		Field []int `env:"TEST_PARSE_SLICE_ErrInvalidElement"`
	}{}
	os.Setenv("TEST_PARSE_SLICE_ErrInvalidElement", "-2,-1,0,x,i")
	if err := envconf.Parse(&cfg); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestParse_Map_Ok(t *testing.T) {
	cfg := struct {
		Field map[int]string `env:"TEST_PARSE_MAP_OK"`
	}{}
	os.Setenv("TEST_PARSE_MAP_OK", "0:ok,1:2,2:test")
	expectedResult := map[int]string{0: "ok", 1: "2", 2: "test"}
	if err := envconf.Parse(&cfg); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Field, expectedResult) {
		t.Fatalf("incorrect result. expected=%v actual=%v", expectedResult, cfg.Field)
	}
}

func TestParse_MapValueInterfaceDefinedValue_Ok(t *testing.T) {
	v1 := &struct {
		Field1 int `default:"5"`
	}{}
	cfg := struct {
		Field map[string]interface{}
	}{
		Field: map[string]interface{}{
			"xyz": v1,
		},
	}

	if err := envconf.Parse(&cfg); err != nil {
		t.Fatal(err)
	}
	if v1.Field1 != 5 {
		t.Fatalf("incorrect result. expected=%v actual=%v", 5, v1.Field1)
	}
}

func TestParse_MapValueInterface_Ok(t *testing.T) {
	cfg := struct {
		Field map[string]interface{} `default:"1:test,2:test2"`
	}{}

	if err := envconf.Parse(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Field["1"] != "test" {
		t.Fatalf("incorrect result. expected=%v actual=%v", "test", cfg.Field["1"])
	}
}

func TestParse_Map_ErrUnsupportedType(t *testing.T) {
	t.Run("InvalidValue", func(t *testing.T) {
		cfg := struct {
			Field map[int]interface{} `env:"TEST_PARSE_MAP_OK_ErrUnsupportedType"`
		}{}
		os.Setenv("TEST_PARSE_MAP_OK_ErrUnsupportedType", "1:1")
		if err := envconf.Parse(&cfg); err != nil {
			t.Fatal("expected nil but got error: ", err)
		}
		if v, ok := cfg.Field[1]; !ok || v != "1" {
			t.Fatalf("unexpected result: %#v", cfg.Field)
		}
	})

	t.Run("InvalidKey", func(t *testing.T) {
		cfg := struct {
			Field map[struct{ k interface{} }]interface{} `env:"TEST_PARSE_MAP_OK_ErrUnsupportedType"`
		}{}
		os.Setenv("TEST_PARSE_MAP_OK_ErrUnsupportedType", "1:1")
		if err := envconf.Parse(&cfg); err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("InvalidKeyValue_KeyInputOnly", func(t *testing.T) {
		cfg := struct {
			Field map[struct{ k interface{} }]interface{} `env:"TEST_PARSE_MAP_OK_ErrUnsupportedType"`
		}{}
		os.Setenv("TEST_PARSE_MAP_OK_ErrUnsupportedType", "1")
		if err := envconf.Parse(&cfg); err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
