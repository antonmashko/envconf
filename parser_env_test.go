package envconf_test

// type testConfigWithPointer struct {
// 	Inner *TestInnerStruct
// }

// type TestInnerStruct struct {
// 	Foo string `env:"ENV_FOO"`
// }

// func TestParseEnvVarOnPointer_Ok(t *testing.T) {
// 	os.Setenv("ENV_FOO", "bar")
// 	var cfg testConfigWithPointer
// 	err := envconf.New().Parse(&cfg)
// 	if err != nil {
// 		t.Fatalf("Parse: %s", err)
// 	}

// 	if cfg.Inner == nil || cfg.Inner.Foo == "" {
// 		t.Fatalf("Incorrect result: %v", cfg.Inner)
// 	}
// }
