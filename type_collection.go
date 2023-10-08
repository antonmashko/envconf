package envconf

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/antonmashko/envconf/external"
	"github.com/antonmashko/envconf/option"
)

type collectionDefiner interface {
	fromInterface(interface{}, option.ConfigSource) (interface{}, error)
	fromString(string, option.ConfigSource) (interface{}, error)
	withoutValue() (interface{}, error)
}

type collectionType struct {
	*configField
	cd  collectionDefiner
	v   reflect.Value
	ext external.ExternalSource
}

func newCollectionType(v reflect.Value, f *configField, cd collectionDefiner) *collectionType {
	return &collectionType{
		configField: f,
		v:           v,
		cd:          cd,
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
	if c.parent() != nil {
		c.ext = external.AsExternalSource(c.Name, c.parent().externalSource())
	}
	var err error
	v, cs := c.Value()
	switch cs {
	case option.NoConfigValue:
		v, err = c.cd.withoutValue()
	case option.ExternalSource:
		v, err = c.cd.fromInterface(v, cs)
	default:
		// value specified for entire collection
		str, ok := v.(string)
		if !ok {
			return ErrUnsupportedType
		}
		v, err = c.cd.fromString(str, cs)
	}

	if err != nil {
		return err
	}
	return c.set(v, cs)
}

type sliceType struct {
	*collectionType
}

func newSliceType(v reflect.Value, f *configField) *sliceType {
	sl := &sliceType{}
	col := newCollectionType(v, f, sl)
	sl.collectionType = col
	return sl
}

func (s *sliceType) fromString(value string, cs option.ConfigSource) (interface{}, error) {
	sl := strings.Split(value, ",")
	switch s.v.Kind() {
	case reflect.Array:
		if len(sl) > s.v.Len() {
			return nil, errors.New("elements in value more than len of array")
		}
	case reflect.Slice:
		s.v.Set(reflect.MakeSlice(s.v.Type(), len(sl), cap(sl)))
	}

	for i := range sl {
		rv := s.v.Index(i)
		st := newDefinedConfigField(sl[i], cs, s,
			reflect.StructField{Name: strconv.Itoa(i), Type: rv.Type()}, s.parser)
		err := s.defineItem(rv, st)
		if err != nil {
			return nil, err
		}
	}
	return sl, nil
}

func (s *sliceType) fromInterface(v interface{}, cs option.ConfigSource) (interface{}, error) {
	if cs != option.ExternalSource {
		return nil, ErrUnsupportedType
	}
	return s.rescan(cs)
}

func (s *sliceType) withoutValue() (interface{}, error) {
	return s.rescan(option.NoConfigValue)
}

func (s *sliceType) rescan(cs option.ConfigSource) (interface{}, error) {
	if !s.v.CanInterface() {
		return nil, errors.New("reflect: cannot interface")
	}
	for i := 0; i < s.v.Len(); i++ {
		rv := s.v.Index(i)
		if !rv.CanInterface() {
			return nil, errors.New("reflect: cannot interface")
		}
		log.Println("TUTA")
		st := newDefinedConfigField(rv.Interface(), cs, s,
			reflect.StructField{Name: strconv.Itoa(i), Type: rv.Type()}, s.parser)
		err := s.defineItem(rv, st)
		if err != nil {
			return nil, err
		}
	}
	return s.v.Interface(), nil
}

type mapType struct {
	*collectionType
}

func newMapType(v reflect.Value, f *configField) *mapType {
	mp := &mapType{}
	col := newCollectionType(v, f, mp)
	mp.collectionType = col
	return mp
}

func (m *mapType) fromString(value string, cs option.ConfigSource) (interface{}, error) {
	sl := strings.Split(value, ",")
	vt := m.v.Type()
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
		st := newDefinedConfigField(value, cs, m,
			reflect.StructField{Name: fmt.Sprint(key), Type: rvalType}, m.parser)
		rvvalue := reflect.New(rvalType).Elem()
		err = m.defineItem(rvvalue, st)
		if err != nil {
			return nil, err
		}
		rmp.SetMapIndex(rvkey, rvvalue)
	}
	m.v.Set(rmp)
	return sl, nil
}

func (m *mapType) fromInterface(v interface{}, cs option.ConfigSource) (interface{}, error) {
	if cs != option.ExternalSource {
		return nil, ErrUnsupportedType
	}
	return m.rescan(cs)
}

func (m *mapType) withoutValue() (interface{}, error) {
	return m.rescan(option.NoConfigValue)
}

func (m *mapType) rescan(cs option.ConfigSource) (interface{}, error) {
	if !m.v.CanInterface() {
		return nil, errors.New("reflect: cannot interface")
	}
	mp := m.v.MapRange()
	for mp.Next() {
		rkey := mp.Key()
		rval := mp.Value()
		if !rkey.CanInterface() {
			return nil, errors.New("reflect: cannot interface map.Key")
		}
		if !rval.CanInterface() {
			return nil, errors.New("reflect: cannot interface map.Value")
		}
		st := newDefinedConfigField(rval.Interface(), cs, m,
			reflect.StructField{Name: fmt.Sprint(rkey.Interface()), Type: rval.Type()}, m.parser)
		err := m.defineItem(rval, st)
		if err != nil {
			return nil, err
		}
	}
	return m.v.Interface(), nil
}
