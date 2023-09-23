package envconf

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// External config source
// Implementation of this interface should be able to `Unmarshal` data into map[string]interface{},
// where interface{} should be also same map type for the nested structures
type External interface {
	// TagName is key name in golang struct tag (json, yaml, toml etc.).
	TagName() string
	// Unmarshal parses the external data and stores the result
	// in the value pointed to by v.
	// Usually, it just wraps the existing `Unmarshal` function of third-party libraries
	Unmarshal(v interface{}) error
}

type emptyExt struct{}

func (emptyExt) TagName() string {
	return ""
}

func (emptyExt) Unmarshal(v interface{}) error { return nil }

type externalConfig struct {
	s    *structType
	ext  External
	data map[string]interface{}
}

func newExternalConfig(ext External) *externalConfig {
	return &externalConfig{
		ext:  ext,
		data: make(map[string]interface{}),
	}
}

func (c *externalConfig) Unmarshal(v interface{}) error {
	if c.ext == (emptyExt{}) {
		return nil
	}
	mp := make(map[string]interface{})
	err := c.ext.Unmarshal(&mp)
	if err != nil {
		return err
	}
	c.data = make(map[string]interface{})
	if err = c.fillMap(c.s, mp); err != nil {
		return err
	}
	return nil
}

func (c *externalConfig) setParentStruct(s *structType) {
	c.s = s
}

func (c *externalConfig) get(f field) (interface{}, bool) {
	v, ok := c.data[fullname(f)]
	return v, ok
}

func (c *externalConfig) fillMap(s *structType, src map[string]interface{}) error {
	for k, v := range src {
		f, ok := c.findField(k, s)
		if !ok {
			continue
		}

		mp, ok := v.(map[string]interface{})
		if ok && f.structField().Type.Kind() != reflect.Map {
			st, ok := f.(*structType)
			if !ok {
				return &Error{
					Message:   fmt.Sprintf("unable to cast %s to struct", f.structField().Type),
					FieldName: fullname(f),
				}
			}
			c.fillMap(st, mp)
			continue
		}

		c.data[fullname(f)] = v
	}
	return nil
}

func (c *externalConfig) findField(key string, s *structType) (field, bool) {
	var fr rune
	for _, r := range key {
		fr = r
		break
	}
	lc := unicode.IsLower(fr)
	for _, f := range s.fields {
		// if annotation exists matching only by it
		sf := f.structField()
		tagName, ok := sf.Tag.Lookup(c.ext.TagName())
		if ok {
			extName := c.validateAndFix(tagName)
			if key == extName {
				return f, true
			}
		}

		// unexportable field. looking for any first match with EqualFold
		if lc && strings.EqualFold(key, sf.Name) {
			return f, true
		}

		if key == sf.Name {
			return f, true
		}
	}

	return nil, false
}

func (c *externalConfig) validateAndFix(name string) string {
	for i, r := range name {
		if r == ',' {
			return name[:i]
		}
	}
	return name
}
