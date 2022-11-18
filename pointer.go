package envconf

import "reflect"

func depointerize(v reflect.Value) reflect.Value {
	// using loop for handling case with more than one referencing
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			tmp := reflect.New(v.Type().Elem())
			v.Set(tmp)
			v = tmp
			continue
		} else {
			v = v.Elem()
		}
	}
	return v
}

