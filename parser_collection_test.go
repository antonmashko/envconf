package envconf_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/antonmashko/envconf"
)

func TestParse_Slice_Ok(t *testing.T) {
	cfg := struct {
		Field []int `env:"TEST_PARSE_SLICE_OK"`
	}{}
	os.Setenv("TEST_PARSE_SLICE_OK", "-2,-1,0,1,2")
	expectedResult := []int{-2, -1, 0, 1, 2}
	if err := envconf.Parse(&cfg); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Field, expectedResult) {
		t.Fatalf("incorrect result. expected=%v actual=%v", expectedResult, cfg.Field)
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

// func TestParse_Map_Ok(t *testing.T) {
// 	cfg := struct {
// 		Field map[int]interface{} `env:"TEST_PARSE_MAP_OK"`
// 	}{}
// 	os.Setenv("TEST_PARSE_MAP_OK", "0:ok,1:2,2:test")
// 	expectedResult := map[int]interface{}{0: "ok", 1: "2", 2: "test"}
// 	if err := envconf.Parse(&cfg); err != nil {
// 		t.Fatal(err)
// 	}
// 	if !reflect.DeepEqual(cfg.Field, expectedResult) {
// 		t.Fatalf("incorrect result. expected=%v actual=%v", expectedResult, cfg.Field)
// 	}
// }
