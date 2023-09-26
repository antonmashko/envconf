package envconf

import (
	"errors"
	"reflect"
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

type collectionSliceType struct {
	*collectionType
}

func (t *collectionSliceType) create(value string) (interface{}, error) {
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
		err := setFromString(t.v.Index(i), sl[i])
		if err != nil {
			return nil, err
		}
	}
	return sl, nil
}

func (t *collectionSliceType) define() error {
	v, p, err := t.readConfigValue()
	if err != nil {
		return err
	}

	if p != option.ExternalSource {
		if value, ok := v.(string); ok {
			v, err = t.create(value)
			if err != nil {
				return err
			}
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

func (t *collectionMapType) create(value string) (interface{}, error) {
	sl := strings.Split(value, ",")
	vt := t.v.Type()
	rmp := reflect.MakeMap(vt)
	key := vt.Key()
	elem := vt.Elem()
	for i := range sl {
		idx := strings.IndexRune(sl[i], ':')
		rvkey := reflect.New(key).Elem()
		rvval := reflect.New(elem).Elem()
		if idx == -1 {
			if err := setFromString(rvkey, sl[i]); err != nil {
				return nil, err
			}
		} else {
			if err := setFromString(rvkey, sl[i][:idx]); err != nil {
				return nil, err
			}
			if err := setFromString(rvval, sl[i][idx+1:]); err != nil {
				return nil, err
			}
		}
		rmp.SetMapIndex(rvkey, rvval)
	}
	t.v.Set(rmp)
	return sl, nil
}

func (t *collectionMapType) define() error {
	v, p, err := t.readConfigValue()
	if err != nil {
		return err
	}

	if p != option.ExternalSource {
		if value, ok := v.(string); ok {
			v, err = t.create(value)
			if err != nil {
				return err
			}
		}
	}

	t.definedValue = &definedValue{
		source: p,
		value:  v,
	}

	return nil
}
