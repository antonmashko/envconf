package envconf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type ErrIncorrectValue struct {
	Value  string
	Symbol string
}

func (e ErrIncorrectValue) Error() string {
	return fmt.Sprintf("%s contains invalid symbol %s\n", e.Value, e.Symbol)
}

var (
	ErrInvalidPair    = errors.New("invalid pair for env variable")
	ErrQuotesInQuotes = errors.New("quotes in quotes")
	ErrNotPairQuotes  = errors.New("not pair quotes")
)

type EnvConfig struct {
	Envs        map[string]string
	IsTrimSpace bool
	Quote       string
}

func NewEnvConf() *EnvConfig {
	return &EnvConfig{Envs: map[string]string{}}
}

func (e *EnvConfig) TrimSpace(t bool) *EnvConfig {
	e.IsTrimSpace = t
	return e
}

func (e *EnvConfig) Parse(data io.Reader) error {
	lines, err := e.readLine(data)
	if err != nil {
		return err
	}
	err = e.parseEnvLines(lines)
	if err != nil {
		return err
	}
	return nil
}

func (e *EnvConfig) parseEnvLines(lines []string) error {
	for _, line := range lines {
		if strings.Contains(line, "#") {
			line = strings.TrimSpace(line)
			line = trimComment(line)
			// continue if block comment
			if len(line) == 0 {
				continue
			}
		}
		line = strings.TrimLeft(line, "export ")
		i := strings.Index(line, "=")
		if i == -1 {
			return ErrInvalidPair
		}
		key := line[:i]
		value := line[i+1:]
		// trim space symbols near key and value
		if e.IsTrimSpace {
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
		} else if strings.HasSuffix(key, " ") || strings.HasPrefix(value, " ") {
			return ErrIncorrectValue{Value: key, Symbol: " "}
		}
		var err error
		key, err = e.trimQuotes(key)
		if err != nil {
			return err
		}
		if strings.Contains(key, " ") {
			return ErrInvalidPair
		}
		value, err = e.trimQuotes(value)
		if err != nil {
			return err
		}
		value = e.trimCharacterEscaping(value)
		e.Envs[key] = value
	}
	return nil
}

func (e *EnvConfig) readLine(data io.Reader) ([]string, error) {
	b := bufio.NewReader(data)
	lines := []string{}
	for {
		l, err := b.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			if len(l) == 0 {
				return lines, nil
			}
		}
		lines = append(lines, l)
	}
}

func (e *EnvConfig) Set() error {
	for k, v := range e.Envs {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (e *EnvConfig) trimQuotes(s string) (string, error) {
	s = strings.TrimSpace(s)
	// string contains only space symbols
	if len(s) == 0 {
		return "", nil
	}
	// valide pair quotes
	if s[0] == '"' || s[0] == '\'' {
		e.Quote = string(s[0])
		if s[0] != s[len(s)-1] {
			return "", ErrNotPairQuotes
		}
		s = s[1 : len(s)-1]
	}
	return s, nil
}

func (e *EnvConfig) trimCharacterEscaping(s string) string {
	if e.Quote == "'" {
		// remove escaping with '
		return escaping(s, "'")
	} else if e.Quote == "\"" {
		// remove escaping with "
		return escaping(s, "\"")
	}
	return s
}

func escaping(s, q string) string {
	for i, j := 0, 1; j < len(s); i, j = i+1, j+1 {
		if string(s[i]) == `\` && string(s[j]) == q {
			s = s[:i] + s[i+1:]
		}
	}
	return s
}

func trimComment(s string) string {
	segmentsBetweenHashes := strings.Split(s, "#")
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
	return strings.Join(segmentsToKeep, "#")
}
