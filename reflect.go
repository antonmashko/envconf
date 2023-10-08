package envconf

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func createFromString(rt reflect.Type, value string) (reflect.Value, interface{}, error) {
	rv := reflect.New(rt)
	v, err := setFromString(rv.Elem(), value)
	if err != nil {
		return reflect.Value{}, nil, err
	}
	return rv.Elem(), v, nil
}

func setFromString(rv reflect.Value, value string) (interface{}, error) {
	oval := value
	value = strings.Trim(value, " ")
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(oval)
		return oval, nil
	case reflect.Bool:
		i, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		rv.SetBool(i)
		return i, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		if _, ok := rv.Interface().(time.Duration); ok {
			d, err := time.ParseDuration(value)
			if err != nil {
				return nil, err
			}
			i = d.Nanoseconds()
		} else {
			var err error
			i, err = strconv.ParseInt(value, 0, rv.Type().Bits())
			if err != nil {
				return nil, err
			}
		}
		rv.SetInt(i)
		return i, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 0, rv.Type().Bits())
		if err != nil {
			return nil, err
		}
		rv.SetUint(i)
		return i, nil
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(value, rv.Type().Bits())
		if err != nil {
			return nil, err
		}
		rv.SetFloat(i)
		return i, nil
	case reflect.Complex64, reflect.Complex128:
		i, err := strconv.ParseComplex(value, rv.Type().Bits())
		if err != nil {
			return nil, err
		}
		rv.SetComplex(i)
		return i, nil
	default:
		return nil, ErrUnsupportedType
	}
}
