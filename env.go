package envconf

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalidPair = errors.New("invalid pair for env variable")
)

type EnvConfig struct{}

func NewEnvConf() *EnvConfig {
	return &EnvConfig{}
}

func (e *EnvConfig) SetEnv(data io.Reader) error {
	r := bufio.NewReader(data)
	for {
		pair, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF && len(pair) == 0 {
			return nil
		}
		i := bytes.Index(pair, []byte("="))
		if i == -1 {
			return ErrInvalidPair
		}
		err = os.Setenv(string(pair[:i]), string(pair[i+1:len(pair)]))
		if err != nil {
			return err
		}
	}
}
