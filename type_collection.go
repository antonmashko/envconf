package envconf

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/antonmashko/envconf/external"
	"github.com/antonmashko/envconf/option"
)

type collectionType struct {
	*fieldType
	parser *EnvConf
	ext    external.ExternalSource
}

func newCollectionType(v reflect.Value, p field, sf reflect.StructField, parser *EnvConf) *collectionType {
	return &collectionType{
		fieldType: newFieldType(v, p, sf, parser.PriorityOrder(), parser.opts.AllowExternalEnvInjection),
		parser:    parser,
		ext:       external.NilContainer{},
	}
}

func (c *collectionType) externalSource() external.ExternalSource {
	return c.ext
}

func (c *collectionType) defineItem(v reflect.Value, p field, sf reflect.StructField, value interface{}, cs option.ConfigSource) error {
	f := createFieldFromValue(v, p, sf, c.parser)
	if err := f.init(); err != nil {
		return err
	}
	c.parser.fieldInitialized(f)
	if err := f.define(); err != nil {
		ft := asFieldType(f)
		if ft != nil && err == ErrConfigurationNotFound {
			err = ft.defineFromValue(value, cs)
			if err != nil {
				c.parser.fieldNotDefined(f, err)
				return err
			}
		} else {
			c.parser.fieldNotDefined(f, err)
			return err
		}
	}
	c.parser.fieldDefined(f)
	return nil
}

type collectionSliceType struct {
	*collectionType
}

func (t *collectionSliceType) createFromString(value string, cs option.ConfigSource) (interface{}, error) {
	sl := strings.Split(value, ",")
	switch t.v.Kind() {
	case reflect.Array:
		if len(sl) > t.v.Len() {
			return nil, errors.New("elements in value more than len of array")
		}
	case reflect.Slice:
		t.v.Set(reflect.MakeSlice(t.v.Type(), len(sl), cap(sl)))
	}

	for i := 0; i < t.v.Len(); i++ {
		rv := t.v.Index(i)
		err := t.defineItem(rv, t, reflect.StructField{Name: strconv.Itoa(i), Type: rv.Type()}, sl[i], cs)
		if err != nil {
			return nil, err
		}
	}
	return sl, nil
}

func (t *collectionSliceType) define() error {
	if t.p != nil {
		t.ext = external.AsExternalSource(t.sf.Name, t.p.externalSource())
	}
	v, p, err := t.readValue()
	if err == nil {
		// value specified for entire collection
		value, ok := v.(string)
		if ok {
			_, err = t.createFromString(value, p)
			if err != nil {
				return err
			}
			return nil
		}
	}

	if !t.v.CanInterface() {
		return nil
	}
	v = t.v.Interface()
	p = option.ConfigSource(-1)

	for i := 0; i < t.v.Len(); i++ {
		rv := t.v.Index(i)
		err = t.defineItem(rv, t, reflect.StructField{Name: strconv.Itoa(i), Type: rv.Type()},
			rv.Interface(), p)
		if err != nil {
			return err
		}
	}

	t.definedValue = &definedValue{
		source: p,
		value:  v,
	}

	return nil
}

type collectionMapType struct {
	*collectionType
}

func (t *collectionMapType) createFromString(value string, p option.ConfigSource) (interface{}, error) {
	sl := strings.Split(value, ",")
	vt := t.v.Type()
	rmp := reflect.MakeMap(vt)
	key := vt.Key()
	elem := vt.Elem()
	for i := range sl {
		idx := strings.IndexRune(sl[i], ':')
		rvkey := reflect.New(key).Elem()
		rvval := reflect.New(elem).Elem()
		var key, value string
		var err error
		if idx == -1 {
			key = sl[i]
		} else {
			key = sl[i][:idx]
			value = sl[i][idx+1:]
		}
		if _, err := setFromString(rvkey, key); err != nil {
			return nil, err
		}
		err = t.defineItem(rvval, t, reflect.StructField{Name: fmt.Sprint(key), Type: rvkey.Type()}, value, p)
		if err != nil {
			return nil, err
		}
		rmp.SetMapIndex(rvkey, rvval)
	}
	t.v.Set(rmp)
	return sl, nil
}

func (t *collectionMapType) define() error {
	if t.p != nil {
		t.ext = external.AsExternalSource(t.sf.Name, t.p.externalSource())
	}
	v, p, err := t.readValue()
	if err == nil {
		// value specified for entire collection
		value, ok := v.(string)
		if ok {
			_, err = t.createFromString(value, p)
			if err != nil {
				return err
			}
			return nil
		}
	}

	if !t.v.CanInterface() {
		return nil
	}
	v = t.v.Interface()
	p = option.ConfigSource(-1)

	mp := t.v.MapRange()
	for mp.Next() {
		rkey := mp.Key()
		rval := mp.Value()
		if !rkey.CanInterface() {
			continue
		}

		var iv interface{}
		if rval.CanInterface() {
			iv = rval.Interface()
		}

		err = t.defineItem(rval, t, reflect.StructField{Name: fmt.Sprint(rkey.Interface()), Type: rval.Type()},
			iv, p)
		if err != nil {
			return err
		}
	}

	t.definedValue = &definedValue{
		source: p,
		value:  v,
	}

	return nil
}
