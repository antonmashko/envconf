package envconf

import (
	"reflect"
	"strings"
)

// External config source
type External interface {
	TagName() string
	Unmarshal(interface{}) error
}

type emptyExt struct{}

func (c *emptyExt) TagName() string {
	return ""
}

func (c *emptyExt) Unmarshal(v interface{}) error { return nil }

type externalConfig struct {
	ext   External
	names map[string]string
	data  map[string]interface{}
}

func newExternalConfig(ext External) *externalConfig {
	return &externalConfig{
		ext:   ext,
		names: make(map[string]string),
		data:  make(map[string]interface{}),
	}
}

func (c *externalConfig) Unmarshal(v interface{}) error {
	mp := make(map[string]interface{})
	err := c.ext.Unmarshal(&mp)
	if err != nil {
		return err
	}
	// normalizing name for the fast access
	c.fillMap("", mp)
	c.names = make(map[string]string) //release resource
	return nil
}

func (c *externalConfig) initName(f field, st reflect.StructField) {
	if f == (emptyField{}) {
		return
	}
	fname := st.Tag.Get(c.ext.TagName())
	if fname == "" {
		fname = f.name()
	}

	if f.parent() != nil && f.parent().name() != "" {
		parentName := fullname(f.parent())
		var pname string
		var ok bool
		for k, v := range c.names {
			if strings.EqualFold(v, parentName) {
				pname = k
				ok = true
				break
			}
		}
		if !ok {
			panic("unable to find parent name: " + fullname(f.parent()))
		}
		fname = pname + fieldNameDelim + fname
	}
	c.names[strings.ToLower(fname)] = fullname(f)
}

func (c *externalConfig) get(f field) (interface{}, bool) {
	v, ok := c.data[fullname(f)]
	return v, ok
}

func (c *externalConfig) fillMap(srcName string, src map[string]interface{}) {
	for k, v := range src {
		if srcName != "" {
			k = srcName + fieldNameDelim + k
		}
		k = strings.ToLower(k)
		mp, ok := v.(map[string]interface{})
		if ok {
			c.fillMap(k, mp)
			continue
		}
		dstName, ok := c.names[k]
		if !ok {
			panic("unable to find dst field name: " + k)
		}
		c.data[dstName] = v
	}
}
