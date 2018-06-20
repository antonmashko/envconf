package envconf

import (
	"testing"
	"time"
)

func TestParsingDurtaionOK(t *testing.T) {
	tc := struct {
		X time.Duration `default:"5m"`
	}{}
	if err := Parse(&tc); err != nil {
		t.Errorf("failed to parse duration. err=%s", err)
	}
	if tc.X.Minutes() != 5.0 {
		t.Errorf("incorrect value was set. %#v", tc.X)
	}
}
