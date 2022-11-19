package envconf

import (
	"reflect"
	"strconv"
)

type primitiveType struct {
	parent   *structType
	v        reflect.Value
	tag      reflect.StructField
	flag     configSource // flag value
	env      configSource // env value
	def      configSource // default value
	required bool         // if it defined true, value should be defined
	desc     string       // description
	isSet    bool
}

func newPrimitiveType(val reflect.Value, parent *structType, tag reflect.StructField) *primitiveType {
	desc := tag.Tag.Get(tagDescription)
	required, _ := strconv.ParseBool(tag.Tag.Get(tagRequired))
	return &primitiveType{
		parent:   parent,
		v:        val,
		tag:      tag,
		flag:     newFlagSource(tag, desc),
		env:      newEnvSource(tag),
		def:      newDefaultValueSource(tag),
		required: required,
		desc:     desc,
		isSet:    false,
	}
}

func (t *primitiveType) Init() error {
	return nil
}

func (t *primitiveType) Define() error {
	// validate reflect value
	if !t.v.IsValid() {
		return errInvalidFiled
	}
	if !t.v.CanSet() {
		return errFiledIsNotSettable
	}

	// create correct parse priority
	priority := priorityOrder()
	for _, p := range priority {
		var source configSource
		switch p {
		case FlagPriority:
			source = t.flag
		case EnvPriority:
			source = t.env
		// case ExternalPriority:
		// 	panic("not supported")
		case DefaultPriority:
			source = t.def
		}

		if str, ok := source.Value(); ok {
			// set value
			t.isSet = true
			return setFromString(t.v, str)
		}
	}

	return errConfigurationNotSpecified
}
