package option

import "testing"

func TestWithCustomExternalInjection(t *testing.T) {
	tcs := []struct {
		name      string
		value     string
		resultVal string
		resultCS  ConfigSource
	}{
		{name: "ValidEnvVar_Ok", value: "${ .env.ENV_VAR }",
			resultVal: "ENV_VAR", resultCS: EnvVariable},
		{name: "InvalidEnvVar_WithoutPrefix_Err", value: ".env.ENV_VAR }",
			resultVal: "", resultCS: NoConfigValue},
		{name: "InvalidEnvVar_WithoutBracketSuffix_Err", value: "${ .env.ENV_VAR",
			resultVal: "", resultCS: NoConfigValue},
		{name: "InvalidEnvVar_WithoutInjectionName_Err", value: "${ .ENV_VAR }",
			resultVal: "", resultCS: NoConfigValue},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			f := WithExternalInjection()
			opts := &Options{}
			f.Apply(opts)
			res, cs := opts.ExternalInjection()(tc.value)
			if res != tc.resultVal || cs != tc.resultCS {
				t.Fatalf("unexpected result: expected={'%s','%s'} actual={'%s','%s'}",
					tc.resultVal, tc.resultCS, res, cs)
			}
		})
	}
}

func TestWithCustomExternalInjection_CustomFunc_Ok(t *testing.T) {
	const suffix = "_test"
	const expectedCS = ExternalSource
	f := WithCustomExternalInjection(func(s string) (string, ConfigSource) {
		return s + suffix, expectedCS
	})
	opts := &Options{}
	f.Apply(opts)
	const value = "abc"
	res, cs := opts.ExternalInjection()(value)
	if res != value+suffix || cs != expectedCS {
		t.Fatalf("unexpected result: expected={'%s','%s'} actual={'%s','%s'}",
			value+suffix, expectedCS, res, cs)
	}
}
