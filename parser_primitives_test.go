package envconf_test

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"testing"
	"time"
	"unsafe"

	"github.com/antonmashko/envconf"
)

func TestPrimitive_ParseFlatStructWithAllPrimitivesFromDefault_Ok(t *testing.T) {
	data := struct {
		Field1  bool    `default:"true"`
		Field2  int     `default:"1"`
		Field3  int8    `default:"2"`
		Field4  int16   `default:"3"`
		Field5  int32   `default:"4"`
		Field6  int64   `default:"5"`
		Field7  uint    `default:"6"`
		Field8  uint8   `default:"7"`
		Field9  uint16  `default:"8"`
		Field10 uint32  `default:"9"`
		Field11 uint64  `default:"10"`
		Field12 float32 `default:"11"`
		Field13 float64 `default:"12"`
		Field14 string  `default:"13"`
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	verify := func(expected interface{}, actual interface{}) {
		if fmt.Sprint(expected) != fmt.Sprint(actual) {
			t.Fatalf("incorrect value. expected=%v actual=%v", expected, actual)
		}
	}
	verify(true, data.Field1)
	for i, v := range []interface{}{
		data.Field2, data.Field3, data.Field4, data.Field5,
		data.Field6, data.Field7, data.Field8, data.Field9,
		data.Field10, data.Field11, data.Field12, data.Field13,
		data.Field14,
	} {
		verify(i+1, v)
	}
}

func TestPrimitive_ParseComplex64_InvalidValueErr(t *testing.T) {
	data1 := struct {
		Field1 complex64 `default:"value1"`
	}{}
	if err := envconf.Parse(&data1); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestParseFlatStructWithAllPrimitivesFromEnv_Ok(t *testing.T) {
	data := struct {
		Field1  bool    `env:"TEST_FIELD_1"`
		Field2  int     `env:"TEST_FIELD_2"`
		Field3  int8    `env:"TEST_FIELD_3"`
		Field4  int16   `env:"TEST_FIELD_4"`
		Field5  int32   `env:"TEST_FIELD_5"`
		Field6  int64   `env:"TEST_FIELD_6"`
		Field7  uint    `env:"TEST_FIELD_7"`
		Field8  uint8   `env:"TEST_FIELD_8"`
		Field9  uint16  `env:"TEST_FIELD_9"`
		Field10 uint32  `env:"TEST_FIELD_10"`
		Field11 uint64  `env:"TEST_FIELD_11"`
		Field12 float32 `env:"TEST_FIELD_12"`
		Field13 float64 `env:"TEST_FIELD_13"`
		Field14 string  `env:"TEST_FIELD_14"`
	}{}

	os.Setenv("TEST_FIELD_1", "1")
	for i := 2; i <= 14; i++ {
		os.Setenv(fmt.Sprint("TEST_FIELD_", i), fmt.Sprint(i-1))
	}

	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	verify := func(expected interface{}, actual interface{}) {
		if fmt.Sprint(expected) != fmt.Sprint(actual) {
			t.Fatalf("incorrect value. expected=%v actual=%v", expected, actual)
		}
	}
	verify(true, data.Field1)
	for i, v := range []interface{}{
		data.Field2, data.Field3, data.Field4, data.Field5,
		data.Field6, data.Field7, data.Field8, data.Field9,
		data.Field10, data.Field11, data.Field12, data.Field13,
		data.Field14,
	} {
		verify(i+1, v)
	}
}

func TestPrimitive_ParseFlatStructWithAllPrimitivesFromFlag_Ok(t *testing.T) {
	os.Args = append(os.Args, "-test-field-1=true")
	for i := 2; i <= 14; i++ {
		os.Args = append(os.Args, fmt.Sprintf("-test-field-%d=%d", i, i-1))
	}
	data := struct {
		Field1  bool    `flag:"test-field-1"`
		Field2  int     `flag:"test-field-2"`
		Field3  int8    `flag:"test-field-3"`
		Field4  int16   `flag:"test-field-4"`
		Field5  int32   `flag:"test-field-5"`
		Field6  int64   `flag:"test-field-6"`
		Field7  uint    `flag:"test-field-7"`
		Field8  uint8   `flag:"test-field-8"`
		Field9  uint16  `flag:"test-field-9"`
		Field10 uint32  `flag:"test-field-10"`
		Field11 uint64  `flag:"test-field-11"`
		Field12 float32 `flag:"test-field-12"`
		Field13 float64 `flag:"test-field-13"`
		Field14 string  `flag:"test-field-14"`
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	verify := func(expected interface{}, actual interface{}) {
		if fmt.Sprint(expected) != fmt.Sprint(actual) {
			t.Fatalf("incorrect value. expected=%v actual=%v", expected, actual)
		}
	}
	verify(true, data.Field1)
	for i, v := range []interface{}{
		data.Field2, data.Field3, data.Field4, data.Field5,
		data.Field6, data.Field7, data.Field8, data.Field9,
		data.Field10, data.Field11, data.Field12, data.Field13,
		data.Field14,
	} {
		verify(i+1, v)
	}
}

func TestPrimitive_ParseBoolFromDefault_InvalidValueErr(t *testing.T) {
	data := struct {
		Field1 bool `default:"test"`
	}{}
	if err := envconf.Parse(&data); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestPrimitive_ParseDurationFromEnv_Ok(t *testing.T) {
	data := struct {
		Dur time.Duration `env:"TEST_DUR"`
	}{}
	os.Setenv("TEST_DUR", "10s")
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	if data.Dur != time.Second*10 {
		t.Fatalf("incorrect value. expected=%s actual=%s", time.Second*10, data.Dur)
	}
}

func TestPrimitive_ParseDurationFromEnv_InvalidValueErr(t *testing.T) {
	data := struct {
		Dur time.Duration `env:"TEST_DUR"`
	}{}
	os.Setenv("TEST_DUR", "test")
	if err := envconf.Parse(&data); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestPrimitive_ParseInt64FromDefault_InvalidValueErr(t *testing.T) {
	data := struct {
		Field1 int64 `default:"test"`
	}{}
	if err := envconf.Parse(&data); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestPrimitive_ParseUint64FromDefault_InvalidValueErr(t *testing.T) {
	data := struct {
		Field1 uint64 `default:"test"`
	}{}
	if err := envconf.Parse(&data); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestPrimitive_ParseFloat64FromDefault_InvalidValueErr(t *testing.T) {
	data := struct {
		Field1 float64 `default:"test"`
	}{}
	if err := envconf.Parse(&data); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestPrimitive_RequiredValue_Ok(t *testing.T) {
	expected := "test123"
	os.Setenv("TEST_FIELD1", expected)
	data := struct {
		Field1 string `env:"TEST_FIELD1" required:"true"`
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal("Parse: ", err)
	}
	if data.Field1 != expected {
		t.Fatalf("incorrect value. expected=%s actual=%s", expected, data.Field1)
	}
}

func TestPrimitive_RequiredValue_Err(t *testing.T) {
	os.Clearenv()
	data := struct {
		Field1 string `env:"TEST_FIELD1" required:"true"`
	}{}
	err := envconf.Parse(&data)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	e, ok := err.(*envconf.Error)
	if !ok {
		t.Fatalf("invalid error type. expected=*envconf.Error actual=%#v", err)
	}
	const expectedFieldName = "Field1"
	if e.FieldName != "Field1" {
		t.Fatalf("incorrect field name in error. expected=%s actual=%s", expectedFieldName, e.FieldName)
	}
}

func TestPrimitive_ParseComplex_Ok(t *testing.T) {
	data := struct {
		Field1 complex64  `default:"64+64i"`
		Field2 complex128 `default:"128+128i"`
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	expected1 := complex64(64 + 64i)
	if data.Field1 != expected1 {
		t.Fatalf("incorrect value. expected=%v actual=%v", data.Field1, expected1)
	}

	expected2 := complex128(128 + 128i)
	if data.Field1 != expected1 {
		t.Fatalf("incorrect value. expected=%v actual=%v", data.Field2, expected2)
	}
}

func TestPrimitive_ParseComplex_InvalidValueErr(t *testing.T) {
	t.Run("Complex64", func(t *testing.T) {
		data1 := struct {
			Field1 complex64 `default:"value1"`
		}{}
		if err := envconf.Parse(&data1); err == nil {
			t.Fatal("expected error but got nil")
		}
	})
	t.Run("Complex128", func(t *testing.T) {
		data1 := struct {
			Field1 complex128 `default:"value1"`
		}{}
		if err := envconf.Parse(&data1); err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}

func TestPrimitive_ParseIPAddress_Ok(t *testing.T) {
	data := struct {
		FieldIPv4 net.IP `default:"192.0.2.1"`
		FieldIPv6 net.IP `default:"2001:db8::68"`
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	expectedIPv4 := net.ParseIP("192.0.2.1")
	if !expectedIPv4.Equal(data.FieldIPv4) {
		t.Fatalf("incorrect value. expected=%v actual=%v", data.FieldIPv4, expectedIPv4)
	}
	expectedIPv6 := net.ParseIP("2001:db8::68")
	if !expectedIPv6.Equal(data.FieldIPv6) {
		t.Fatalf("incorrect value. expected=%v actual=%v", data.FieldIPv4, expectedIPv6)
	}
}

func TestPrimitive_ParseURL_Ok(t *testing.T) {
	data := struct {
		URL1 *url.URL `default:"http://test.com/"`
		URL2 url.URL  `default:"http://test.com/"`
	}{}
	if err := envconf.Parse(&data); err != nil {
		t.Fatal(err)
	}
	expectedURL, err := url.Parse("http://test.com/")
	if err != nil {
		t.Fatal(err)
	}
	if expectedURL.String() != data.URL1.String() || expectedURL.String() != (&data.URL2).String() {
		t.Fatalf("incorrect value. expected=%v actual1=%v actual2=%v",
			expectedURL.String(), data.URL1.String(), data.URL2.String())
	}
}

func TestPrimitive_ParseURL_ErrInvalidURL(t *testing.T) {
	data := struct {
		URL1 *url.URL `default:";test:test"`
	}{}
	if err := envconf.Parse(&data); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestParse_UnsupportedFieldWithoutPanic_Ok(t *testing.T) {
	cfg := struct {
		Field1 [10]int
		Field2 map[string]interface{}
		Field3 chan struct{}
		Field4 func()
		Field5 unsafe.Pointer
		Field6 uintptr
		Field7 interface{}
	}{}

	if err := envconf.Parse(&cfg); err != nil {
		t.Fatal(err)
	}
}
