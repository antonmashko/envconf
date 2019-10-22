package envconf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

var (
	ErrInvalidPair = errors.New("invalid pair for env variable")
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
		n, err := b.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			if len(n) == 0 {
				return nil
			}
		}
		// trim comment, but save # if the string is in quotation marks
		n = trimComment(n)
		if len(n) == 0 {
			return nil
		}
		fmt.Println(n)
		// Split env variable to key and value
		i := strings.Index(n, "=")
		if i == -1 {
			return ErrInvalidPair
		}
		key, value := strings.TrimFunc(n[0:i], isTrim), strings.TrimFunc(n[i+1:], isTrim)
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

func isTrim(r rune) bool {
	return unicode.IsSpace(r) || r == '"' || r == '\'' || r == '`'
}

func trimComment(line string) string {
	if strings.Contains(line, "#") {
		segmentsBetweenHashes := strings.Split(line, "#")
		quotesAreOpen := false
		var segmentsToKeep []string
		for _, segment := range segmentsBetweenHashes {
			if strings.Count(segment, "\"") == 1 || strings.Count(segment, "'") == 1 {
				if quotesAreOpen {
					quotesAreOpen = false
					segmentsToKeep = append(segmentsToKeep, segment)
				} else {
					quotesAreOpen = true
				}
			}
			if len(segmentsToKeep) == 0 || quotesAreOpen {
				segmentsToKeep = append(segmentsToKeep, segment)
			}
		}
		line = strings.Join(segmentsToKeep, "#")
	}
	return line
}
