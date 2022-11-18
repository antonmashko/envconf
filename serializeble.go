package envconf

import "reflect"

type Serializable interface {
	Serialize(reflect.Value, string) error
}
