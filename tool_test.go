package envconf_test

import (
	"testing"
)

func Fatal(t *testing.T, expected interface{}, actual interface{}, msg string) {
	header := "unexpected result"
	if msg != "" {
		header += " " + msg
	}
	t.Fatalf("%s. expected='%s' actual='%s'", header, expected, actual)
}
