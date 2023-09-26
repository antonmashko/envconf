package envconf

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/antonmashko/envconf/option"
)

type collectionType struct {
	*fieldType
	parser *EnvConf
}

func newCollectionType(v reflect.Value, p field, sf reflect.StructField, parser *EnvConf) *collectionType {
	return &collectionType{
		fieldType: newFieldType(v, p, sf, parser.external, parser.PriorityOrder()),
		parser:    parser,
	}
}

func (c *collectionType) onDefine(v reflect.Value, p field, sf reflect.StructField, value interface{}, cs option.ConfigSource) error {
	f := createFieldFromValue(v, p, sf, c.parser)
	if ft, ok := f.(*fieldType); ok {
		ft.definedValue = &definedValue{
			source: cs,
			value:  value,
		}
	}
	if err := f.init(); err != nil {
		return err
	}
	c.parser.fieldInitialized(f)
	if err := f.define(); err != nil {
		c.parser.fieldNotDefined(f, err)
		return err
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
		v, err := setFromString(rv, sl[i])
		if err != nil {
			return nil, err
		}
		err = t.onDefine(rv, t, reflect.StructField{Name: strconv.Itoa(i), Type: rv.Type()}, v, cs)
		if err != nil {
			return nil, err
		}
	}
	return sl, nil
}

func (t *collectionSliceType) define() error {
	v, p, err := t.readConfigValue()
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
		if !rv.CanInterface() {
			continue
		}
		err = t.onDefine(rv, t, reflect.StructField{Name: strconv.Itoa(i), Type: rv.Type()},
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
		var key, value interface{}
		var err error
		if idx == -1 {
			if key, err = setFromString(rvkey, sl[i]); err != nil {
				return nil, err
			}
		} else {
			if key, err = setFromString(rvkey, sl[i][:idx]); err != nil {
				return nil, err
			}
			if value, err = setFromString(rvval, sl[i][idx+1:]); err != nil {
				return nil, err
			}
		}
		rmp.SetMapIndex(rvkey, rvval)
		err = t.onDefine(rvval, t, reflect.StructField{Name: fmt.Sprint(key), Type: rvkey.Type()}, value, p)
		if err != nil {
			return nil, err
		}
	}
	t.v.Set(rmp)
	return sl, nil
}

func (t *collectionMapType) define() error {
	v, p, err := t.readConfigValue()
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

		err = t.onDefine(rval, t, reflect.StructField{Name: fmt.Sprint(rkey.Interface()), Type: rval.Type()},
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
