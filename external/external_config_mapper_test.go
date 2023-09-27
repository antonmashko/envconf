package external

import (
	"strconv"
	"testing"

	"github.com/antonmashko/envconf/external/json"
)

func TestNilContainer_Ok(t *testing.T) {
	nc := NilContainer{}
	v, ok := nc.Read("")
	if v != nil || ok {
		t.Fatal("unexpected result")
	}
}

func TestSliceContainer_InvalidIdx_Err(t *testing.T) {
	nc := sliceContainer{}
	v, ok := nc.Read("")
	if v != nil || ok {
		t.Fatal("unexpected result")
	}
}

func TestSliceContainer_IdxGreaterThanLen_Err(t *testing.T) {
	nc := sliceContainer{
		1, 2, 3,
	}
	v, ok := nc.Read(strconv.Itoa(len(nc) + 1))
	if v != nil || ok {
		t.Fatal("unexpected result")
	}
}

func TestAsExternalSource_Nil_Err(t *testing.T) {
	es := AsExternalSource("test", nil)
	if es != (NilContainer{}) {
		t.Fatal("unexpected result")
	}
}

func TestAsExternalSource_NilContainer_Err(t *testing.T) {
	es := AsExternalSource("test", NilContainer{})
	if es != (NilContainer{}) {
		t.Fatal("unexpected result")
	}
}

func TestAsExternalSource_NotContainerType_Err(t *testing.T) {
	es := AsExternalSource("foo", mapContainer{"foo": "bar"})
	if es != (NilContainer{}) {
		t.Fatal("unexpected result")
	}
}

func TestExternalConfigMapper_NilExternal_Ok(t *testing.T) {
	extMp := NewExternalConfigMapper(nil)
	if err := extMp.Unmarshal(nil); err != nil {
		t.Fatalf("expected nil but got error. %s", err)
	}
	if extMp.Data() == nil {
		t.Fatal("mapper.Data() is nil")
	}
}

func TestExternalConfigMapper_InvalidJson_Ok(t *testing.T) {
	extMp := NewExternalConfigMapper(json.Json([]byte("<test></test>")))
	result := struct{}{}
	if err := extMp.Unmarshal(&result); err == nil {
		t.Fatal("expected error but got nil")
	}
}
