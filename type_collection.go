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

var _ field = (*collectionSliceType)(nil)
var _ field = (*collectionMapType)(nil)

type collectionType struct {
	*configField
	v   reflect.Value
	p   field
	ext external.ExternalSource
}

func newCollectionType(v reflect.Value, f *configField) *collectionType {
	return &collectionType{
		configField: f,
		v:           v,
		ext:         external.NilContainer{},
	}
}

func (c *collectionType) externalSource() external.ExternalSource {
	return c.ext
}

func (c *collectionType) defineItem(v reflect.Value, cf *configField) error {
	f := createFieldFromValue(v, cf)
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

func (c *collectionType) init() error {
	return c.configField.init(c)
}

func (c *collectionType) define() error {
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

	for i := range sl {
		rv := t.v.Index(i)
		st := newDefinedConfigField(sl[i], cs, t,
			reflect.StructField{Name: strconv.Itoa(i), Type: rv.Type()}, t.parser)
		err := t.defineItem(rv, st)
		if err != nil {
			return nil, err
		}
	}
	return sl, nil
}

func (t *collectionSliceType) define() error {
	if t.p != nil {
		t.ext = external.AsExternalSource(t.Name, t.p.externalSource())
	}
	v, cs := t.Value()
	if cs != option.NoConfigValue {
		// return ErrConfigurationNotFound
		if cs != option.ExternalSource {
			// value specified for entire collection
			value, ok := v.(string)
			if ok {
				var err error
				v, err = t.createFromString(value, cs)
				if err != nil {
					return err
				}
				return t.set(v, cs)
			}
			// return ErrUnsupportedType
		}
	}

	if !t.v.CanInterface() {
		return errors.New("reflect: cannot interface")
	}
	v = t.v.Interface()
	for i := 0; i < t.v.Len(); i++ {
		rv := t.v.Index(i)
		if !rv.CanInterface() {
			return errors.New("reflect: cannot interface")
		}
		st := newDefinedConfigField(rv.Interface(), cs, t,
			reflect.StructField{Name: strconv.Itoa(i), Type: rv.Type()}, t.parser)
		err := t.defineItem(rv, st)
		if err != nil {
			return err
		}
	}
	return t.set(v, cs)
}

type collectionMapType struct {
	*collectionType
}

func (t *collectionMapType) createFromString(value string, cs option.ConfigSource) (interface{}, error) {
	sl := strings.Split(value, ",")
	vt := t.v.Type()
	rmp := reflect.MakeMap(vt)
	rkeyType := vt.Key()
	rvalType := vt.Elem()
	for i := range sl {
		idx := strings.IndexRune(sl[i], ':')
		var key, value string
		if idx == -1 {
			key = sl[i]
		} else {
			key = sl[i][:idx]
			value = sl[i][idx+1:]
		}
		rvkey, _, err := createFromString(rkeyType, key)
		if err != nil {
			return nil, err
		}
		st := newDefinedConfigField(value, cs, t,
			reflect.StructField{Name: fmt.Sprint(key), Type: rvalType}, t.parser)
		rvvalue := reflect.New(rvalType).Elem()
		err = t.defineItem(rvvalue, st)
		if err != nil {
			return nil, err
		}
		rmp.SetMapIndex(rvkey, rvvalue)
	}
	t.v.Set(rmp)
	return sl, nil
}

func (t *collectionMapType) define() error {
	if t.p != nil {
		t.ext = external.AsExternalSource(t.Name, t.p.externalSource())
	}
	v, cs := t.Value()
	if cs != option.NoConfigValue {
		// return ErrConfigurationNotFound
		if cs != option.ExternalSource {
			// value specified for entire collection
			value, ok := v.(string)
			if ok {
				var err error
				v, err = t.createFromString(value, cs)
				if err != nil {
					return err
				}
				return t.set(v, cs)
			}
			// return ErrUnsupportedType
		}
	}

	if !t.v.CanInterface() {
		return nil
	}
	v = t.v.Interface()
	mp := t.v.MapRange()
	for mp.Next() {
		rkey := mp.Key()
		rval := mp.Value()
		if !rkey.CanInterface() {
			return errors.New("reflect: cannot interface map.Key")
		}
		if !rval.CanInterface() {
			return errors.New("reflect: cannot interface map.Value")
		}
		st := newDefinedConfigField(rval.Interface(), cs, t,
			reflect.StructField{Name: fmt.Sprint(rkey.Interface()), Type: rval.Type()}, t.parser)
		err := t.defineItem(rval, st)
		if err != nil {
			return err
		}
	}
	return t.set(v, cs)
}
