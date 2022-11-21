package envconf

// type TestInnerStruct struct {
// 	Data string `env:"DATA"`
// }

// type testStruct struct {
// 	// Str   fmt.Stringer
// 	Inner ***TestInnerStruct `prefix:"INNER_"`
// }

// func TestStructSerialization_Ok(t *testing.T) {
// 	os.Setenv("DATA", "foo")
// 	var input testStruct
// 	err := New().Parse(&input)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if (**input.Inner).Data != "foo" {
// 		t.Fatal("no config")
// 	}
// }
