package external

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type ExternalSource interface {
	Read(string) (interface{}, bool)
}

type NilContainer struct{}

func (NilContainer) Read(index string) (interface{}, bool) {
	return nil, false
}

type sliceContainer []interface{}

func (c sliceContainer) Read(index string) (interface{}, bool) {
	idx, err := strconv.Atoi(index)
	if err != nil {
		return nil, false
	}
	if idx >= len(c) {
		return nil, false
	}
	return c[idx], true
}

type mapContainer map[string]interface{}

func (c mapContainer) Read(key string) (interface{}, bool) {
	ic, ok := c[key]
	return ic, ok
}

func AsExternalSource(name string, es ExternalSource) ExternalSource {
	if es == nil {
		return NilContainer{}
	}
	v, ok := es.Read(name)
	if !ok {
		return NilContainer{}
	}
	switch vt := v.(type) {
	case map[string]interface{}:
		return mapContainer(vt)
	case []interface{}:
		return sliceContainer(vt)
	default:
		return NilContainer{}
	}
}

type ExternalConfigMapper struct {
	ext  External
	data map[string]interface{}
}

func NewExternalConfigMapper(ext External) *ExternalConfigMapper {
	return &ExternalConfigMapper{
		ext:  ext,
		data: map[string]interface{}{},
	}
}

func (c *ExternalConfigMapper) Data() ExternalSource {
	return mapContainer(c.data)
}

func (c *ExternalConfigMapper) Unmarshal(v interface{}) error {
	if c.ext == nil {
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

func (c *ExternalConfigMapper) normalizeMap(rv reflect.Value, mp map[string]interface{}) (map[string]interface{}, error) {
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

func (c *ExternalConfigMapper) normalizeSlice(rv reflect.Value, sl []interface{}) ([]interface{}, error) {
	for i := range sl {
		v, err := c.normalize(rv.Index(i), sl[i])
		if err != nil {
			return nil, err
		}
		sl[i] = v
	}
	return sl, nil
}

func (c *ExternalConfigMapper) normalize(rv reflect.Value, v interface{}) (interface{}, error) {
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
			return nil, fmt.Errorf("unable to cast map[string]interface{} into %s", rv.Type().Name())
		}
	case []interface{}:
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			return c.normalizeSlice(rv, vt)
		default:
			return nil, fmt.Errorf("unable to cast []interface{} into %s", rv.Type().String())
		}
	default:
		return vt, nil
	}
}

func (c *ExternalConfigMapper) equal(key string, lc bool, sf reflect.StructField) bool {
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
