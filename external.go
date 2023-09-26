package envconf

import (
	"fmt"
	"reflect"
	"strconv"
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

func (c *externalConfig) unmarshal(v interface{}) error {
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
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	c.data, err = c.normalizeMap(rv, mp)
	if err != nil {
		return err
	}
	return nil
}

func (c *externalConfig) readFrom(key string, ic interface{}) (interface{}, bool) {
	switch vt := ic.(type) {
	case map[string]interface{}:
		var ok bool
		ic, ok = vt[key]
		return ic, ok
	case []interface{}:
		idx, err := strconv.Atoi(key)
		if err != nil {
			return nil, false
		}
		if idx >= len(vt) {
			return nil, false
		}
		return vt[idx], true
	default:
		return nil, false
	}
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
	var ic interface{} = c.data
	var ok bool
	for i := 0; i < len(path); i++ {
		ic, ok = c.readFrom(path[i].Name, ic)
		if !ok {
			return nil, false
		}
		if ok && i == len(path)-1 {
			return ic, true
		}
	}
	return nil, false
}

func (c *externalConfig) normalizeMap(rv reflect.Value, mp map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for k, v := range mp {
		var fr rune
		for _, r := range k {
			fr = r
			break
		}
		lc := unicode.IsLower(fr)
		// normalizing names(keys) in map
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			sf := rt.Field(i)
			f := rv.Field(i)
			if !c.equal(k, lc, sf) {
				continue
			}
			val, err := c.normalize(f, v)
			if err != nil {
				return nil, err
			}
			result[sf.Name] = val
			break
		}
	}
	return result, nil
}

func (c *externalConfig) normalizeSlice(rv reflect.Value, sl []interface{}) ([]interface{}, error) {
	for i := range sl {
		v, err := c.normalize(rv.Index(i), sl[i])
		if err != nil {
			return nil, err
		}
		sl[i] = v
	}
	return sl, nil
}

func (c *externalConfig) normalize(rv reflect.Value, v interface{}) (interface{}, error) {
	switch vt := v.(type) {
	case map[string]interface{}:
		switch rv.Kind() {
		case reflect.Map:
			return vt, nil
		case reflect.Struct:
			return c.normalizeMap(rv, vt)
		case reflect.Interface:
			if rv.IsValid() && !rv.IsZero() {
				return c.normalize(rv.Elem(), v)
			}
			return vt, nil
		case reflect.Pointer:
			if rv.IsValid() && !rv.IsZero() {
				return c.normalize(rv.Elem(), v)
			}
			return vt, nil
		default:
			return nil, &Error{
				Message: fmt.Sprint("unable to cast map[string]interface{} into ", rv.Type().Name()),
			}
		}
	case []interface{}:
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			return c.normalizeSlice(rv, vt)
		default:
			return nil, &Error{
				Message: fmt.Sprint("unable to cast []interface{} into ", rv.Type().String()),
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
