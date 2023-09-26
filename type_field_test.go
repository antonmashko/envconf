package envconf

import (
	"reflect"
	"testing"
)

func Test_setFromString(t *testing.T) {
	type args struct {
		field reflect.Value
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := setFromString(tt.args.field, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("setFromString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
