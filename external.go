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
	TagName() []string
	// Unmarshal parses the external data and stores the result
	// in the value pointed to by v.
	// Usually, it just wraps the existing `Unmarshal` function of third-party libraries
	Unmarshal(v interface{}) error
}

type emptyExt struct{}

func (emptyExt) TagName() []string {
	return []string{}
}

func (emptyExt) Unmarshal(v interface{}) error { return nil }

type externalConfig struct {
	ext  External
	data map[string]interface{}
}

func newExternalConfig(ext External) *externalConfig {
	return &externalConfig{
		ext:  ext,
		data: make(map[string]interface{}),
	}
}

func (c *externalConfig) unmarshal(rf reflect.Type, v interface{}) error {
	if c.ext == (emptyExt{}) {
		return nil
	}
	err := c.ext.Unmarshal(v)
	if err != nil {
		return err
	}
	mp := make(map[string]interface{})
	err = c.ext.Unmarshal(&mp)
	if err != nil {
		return err
	}
	c.data, err = c.normalizeMap(rf, mp)
	if err != nil {
		return err
	}
	return nil
}

func (c *externalConfig) get(f field) (interface{}, bool) {
	if c.ext == (emptyExt{}) {
		return nil, false
	}
	path := []reflect.StructField{f.structField()}
	for f.parent() != nil {
		path = append([]reflect.StructField{f.parent().structField()}, path...)
		f = f.parent()
	}

	// ignoring top level struct
	path = path[1:]
	var mp map[string]interface{} = c.data
	if len(path) > 1 {
		for i := 0; i < len(path)-1; i++ {
			mp = mp[path[i].Name].(map[string]interface{})
		}
	}

	v, ok := mp[path[len(path)-1].Name]
	return v, ok
}

func (c *externalConfig) normalizeMap(rt reflect.Type, mp map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for k, v := range mp {
		var fr rune
		for _, r := range k {
			fr = r
			break
		}
		lc := unicode.IsLower(fr)
		// normalizing names(keys) in map
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if !c.equal(k, lc, f) {
				continue
			}
			val, err := c.normalize(f.Type, v)
			if err != nil {
				return nil, err
			}
			result[f.Name] = val
			break
		}
	}
	return result, nil
}

func (c *externalConfig) normalizeSlice(rt reflect.Type, sl []interface{}) ([]interface{}, error) {
	for i := range sl {
		v, err := c.normalize(rt, sl[i])
		if err != nil {
			return nil, err
		}
		sl[i] = v
	}
	return sl, nil
}

func (c *externalConfig) normalize(rf reflect.Type, v interface{}) (interface{}, error) {
	switch vt := v.(type) {
	case map[string]interface{}:
		switch rf.Kind() {
		case reflect.Map:
			return vt, nil
		case reflect.Struct:
			return c.normalizeMap(rf, vt)
		default:
			return nil, &Error{
				Message: fmt.Sprint("unable to cast map[string]interface{} into ", rf.String()),
			}
		}
	case []interface{}:
		switch rf.Kind() {
		case reflect.Slice, reflect.Array:
			return c.normalizeSlice(rf.Elem(), vt)
		default:
			return nil, &Error{
				Message: fmt.Sprint("unable to cast []interface{} into ", rf.String()),
			}
		}
	default:
		return vt, nil
	}
}

func (c *externalConfig) equal(key string, lc bool, sf reflect.StructField) bool {
	for _, tagName := range c.ext.TagName() {
		tag, ok := sf.Tag.Lookup(tagName)
		if ok {
			idx := strings.IndexRune(tag, ',')
			if idx != -1 {
				tag = tag[:idx]
			}
			if key == tag {
				return true
			}
		}
	}

	// unexportable field. looking for any first match with EqualFold
	if lc && strings.EqualFold(key, sf.Name) {
		return true
	}

	if key == sf.Name {
		return true
	}

	return false
}
