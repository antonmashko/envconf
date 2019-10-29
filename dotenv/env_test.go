package envconf

import (
	"os"
	"testing"
)

func TestEnvConfigStatusOK(t *testing.T) {
	type testCase struct {
		filePath       string
		expectedValues map[string]string
		isTrimSpace    bool
	}
	testcases := []testCase{
		{
			filePath: "./fixtures/basic.env",
			expectedValues: map[string]string{
				"OPTION_A":     "foo",
				"OPTION_B":     "foo_bar",
				"OPTION_C":     "3.14",
				"OPTION_D":     "42",
				"_OPTION_E":    "foo",
				"____OPTION_F": "bar",
				"OPTION_G":     "1",
				"OPTION_H":     "2",
			},
		},
		{
			filePath: "./fixtures/exported.env",
			expectedValues: map[string]string{
				"OPTION_A": "2",
				"OPTION_B": `\n`,
			},
		},
		{
			filePath: "./fixtures/space.env",
			expectedValues: map[string]string{
				"OPTION_A": "1",
				"OPTION_B": "2",
				"OPTION_C": "3",
				"OPTION_D": "4",
				"OPTION_E": "5",
				"OPTION_F": "",
				"OPTION_G": "",
			},
			isTrimSpace: true,
		},
		{
			filePath: "./fixtures/quoted.env",
			expectedValues: map[string]string{
				"OPTION_A": "1",
				"OPTION_B": "2",
				"OPTION_C": "",
				"OPTION_D": `\n`,
				"OPTION_E": "1",
				"OPTION_F": "2",
				"OPTION_G": "",
				"OPTION_H": `\n`,
				"OPTION_I": `foo 'bar'`,
				"OPTION_J": `foo"bar"`,
				"OPTION_K": `"foo`,
				"OPTION_L": `foo "bar"`,
				"OPTION_M": `foo \bar\`,
				"OPTION_N": `\\foo`,
				"OPTION_O": `foo \"bar\"`,
				"OPTION_P": "`foo bar`",
			},
		},
		{
			filePath: "./fixtures/comment.env",
			expectedValues: map[string]string{
				"OPTION_A": "1",
				"OPTION_B": "2",
			},
		},
	}
	for _, tc := range testcases {
		envsFile, err := os.Open(tc.filePath)
		if err != nil {
			t.Error("cannot open file via path:", tc.filePath)
		}
		envf := NewEnvConf().TrimSpace(tc.isTrimSpace)
		if err := envf.Parse(envsFile); err != nil {
			t.Errorf("cannot read file %s error: %s", tc.filePath, err)
		}
		// comparing result envs with expected
		if len(tc.expectedValues) != len(envf.Envs) {
			t.Errorf("he expected value is not equal to the value from the file: in the test case=%v != in the env file=%v", len(tc.expectedValues), len(envf.Envs))
		}
		for k, v := range envf.Envs {
			values, ok := tc.expectedValues[k]
			if !ok {
				t.Error("expected values not contains key:", k)
			}
			if values != v {
				t.Errorf("the expected value is not equal to the value from the file: in the test case=%v !=in the file=%v. Test key: %s", values, v, k)
			}
		}
	}
}

func TestEnvConfigParseIncorrectFileStatusError(t *testing.T) {
	type testCase struct {
		filePath    string
		isTrimSpace bool
	}
	testcases := []testCase{
		{
			filePath: "./fixtures/space.env",
		},
	}
	for _, tc := range testcases {
		envsFile, err := os.Open(tc.filePath)
		if err != nil {
			t.Error("cannot open file via path:", tc.filePath)
		}
		envf := NewEnvConf().TrimSpace(tc.isTrimSpace)
		err = envf.Parse(envsFile)
		err, ok := err.(ErrIncorrectValue)
		if !ok {
			t.Errorf("incorrect error type: %T", err)
		}
	}
}

func TestEnvConfigParseIncorrectKeyStatusError(t *testing.T) {
	type testCase struct {
		filePath string
	}
	testcases := []testCase{
		{
			filePath: "./fixtures/nagative_export.env",
		},
	}
	for _, tc := range testcases {
		envsFile, err := os.Open(tc.filePath)
		if err != nil {
			t.Error("cannot open file via path:", tc.filePath)
		}
		envf := NewEnvConf()
		err = envf.Parse(envsFile)
		if err != ErrInvalidPair {
			t.Errorf("incorrect error type: %T", err)
		}
	}
}
