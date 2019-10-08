package envconf

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
)

var (
	ErrInvalidPair = errors.New("invalid pair for env variable")
)

const (
	commentSymbol = '#'
)

type EnvConfig struct {
	envs map[string]string
}

func NewEnvConf() *EnvConfig {
	return &EnvConfig{envs: map[string]string{}}
}

func (e *EnvConfig) Parse(data io.Reader) error {
	b := bufio.NewReader(data)
	for {
		n, err := b.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			if len(n) == 0 {
				return nil
			}
		}
		n = bytes.Trim(n, " ")
		n = bytes.Trim(n, "\n")
		n = bytes.Trim(n, "\t")
		if n[0] == commentSymbol {
			continue
		}
		i := bytes.Index(n, []byte("="))
		if i == -1 {
			return ErrInvalidPair
		}
		key := string(n[:i])
		value := string(n[i+1 : len(n)])
		if isEmptyLine(value) {
			e.envs[key] = ""
			continue
		}
		if isQuotes(value, 0, len(value)-1) {
			value = cutQuotes(value)
		}
		e.envs[key] = value
	}
}

func (e *EnvConfig) Set() error {
	for k, v := range e.envs {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func isEmptyLine(s string) bool {
	return isQuotes(s, 0, 1)
}

func cutQuotes(s string) string {
	if isQuotes(s, 0, len(s)-1) {
		if string(s[0]) == `"` {
			s = strings.Trim(s, `"`)
		} else {
			s = strings.Trim(s, `'`)
		}
	}
	return s
}

func isQuotes(s string, first, second int) bool {
	open := string(s[first])
	close := string(s[second])
	if (open == `"` || open == `'`) && open == close {
		return true
	}
	return false
}
