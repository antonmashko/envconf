package envconf

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestSimpleExternalEnvConfigOK(t *testing.T) {
	env := bytes.NewBuffer([]byte("Name=illia"))
	tc := struct {
		Name string `env:"name"`
	}{}
	NewEnvConf().SetEnv(env)
	err := Parse(&tc)
	if err != nil {
		t.Error(err)
	}
	if tc.Name != "illia" {
		t.Error("not equatl:", tc.Name)
	}
	fmt.Println(os.Getenv("Name"))
}
