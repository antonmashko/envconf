package envconf

import "reflect"

var ErrUnsupportedType = errUnsupportedType

type structData struct {
	name   string
	parent *structData
	values map[string]Value

	val reflect.Value
	tag reflect.StructTag
}

func (s *structData) Serialize(v reflect.Value, data string) error {
	return s.handleStruct(v)
}

func (s *structData) handleStruct(v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		var err error
		switch v.Kind() {
		case reflect.Struct:
			err = s.Serialize(v, "")
		default:
			newValue(nil, v, v.Type().Field(i))
		}
		if err != nil {
			return err
		}
	}
	return nil
}
