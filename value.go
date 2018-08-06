package envconf

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	tagFlag        = "flag"
	tagEnv         = "env"
	tagDefault     = "default"
	tagRequired    = "required"
	tagDescription = "description"
	tagIgnored     = "-"
	tagNotDefined  = ""

	valIgnored    = "ignored"
	valNotDefined = "N/D"
	valDefault    = "*"
)

var (
	//errors
	errInvalidFiled       = errors.New("invalid field")
	errFiledIsNotSettable = errors.New("field is not settable")
	errUnsupportedType    = errors.New("unsupported type")
	errRequiredFiled      = errors.New("required field")
)

type flagV struct {
	name    string
	v       string
	defined bool
}

func (f *flagV) define(tag reflect.StructField, usage string) {
	f.name = tag.Tag.Get(tagFlag)
	if f.name == tagNotDefined {
		f.name = tagIgnored
	} else if strings.ToLower(f.name) == valDefault {
		f.name = strings.ToLower(tag.Name)
	}
	if f.name != tagIgnored {
		flag.Var(f, f.name, usage)
	}
}

func (f *flagV) value() (string, bool) {
	if f.name == tagIgnored {
		return "", false
	}
	return f.v, f.defined
}

func (f *flagV) Set(value string) error {
	f.v = value
	f.defined = true
	return nil
}

func (f *flagV) String() string {
	return f.v
}

type envV struct {
	name string
}

func (e *envV) define(tag reflect.StructField) {
	//define env name
	e.name = tag.Tag.Get(tagEnv)
	if e.name == tagNotDefined {
		e.name = tagIgnored
	} else if e.name == valDefault {
		e.name = strings.ToUpper(tag.Name)
	}
}

func (e *envV) value() (string, bool) {
	if e.name != tagIgnored {
		return os.LookupEnv(e.name)
	}
	return "", false
}

type defaultV struct {
	defined bool
	v       string
}

func (e *defaultV) define(tag reflect.StructField) {
	e.v, e.defined = tag.Tag.Lookup(tagDefault)
}

func (e *defaultV) value() (string, bool) {
	return e.v, e.defined
}

type Value interface {
	Owner() Value
	Name() string
	Tag() reflect.StructField
}

type value struct {
	owner    *parser
	field    reflect.Value
	tag      reflect.StructField
	flagV    flagV    // flag value
	envV     envV     // env value
	defaultV defaultV // default value
	required bool     // if it defined true, value should be defined
	desc     string   // description
}

func newValue(owner *parser, field reflect.Value, tag reflect.StructField) *value {
	v := &value{field: field, tag: tag, owner: owner}
	// Parse description
	v.desc = tag.Tag.Get(tagDescription)
	(&v.flagV).define(tag, v.desc)
	(&v.envV).define(tag)
	(&v.defaultV).define(tag)
	// Parse required
	rq := tag.Tag.Get(tagRequired)
	v.required, _ = strconv.ParseBool(rq)
	return v
}

func (v *value) Owner() Value {
	if v.owner == nil {
		return nil
	}
	return v.owner
}

func (v *value) Name() string {
	return v.tag.Name
}

func (v *value) Tag() reflect.StructField {
	return v.tag
}

func (v *value) fullname() string {
	result := v.Name()
	owner := v.owner
	for owner != nil && owner.Name() != "" {
		result = owner.Name() + string(Separator) + result
		owner = owner.parent
	}
	return result
}

// func (v *value) name() string {
// 	op := v.owner.Path()
// 	if op != "" {
// 		op += string(Separator)
// 	}
// 	return op + v.tag.Name
// }

func (v *value) define() error {
	ferr := func(err error) error {
		if v.required {
			return err
		}
		return nil
	}
	// validate reflect value
	if !v.field.IsValid() {
		return ferr(errInvalidFiled)
	}
	if !v.field.CanSet() {
		return ferr(errFiledIsNotSettable)
	}
	if v.field.Kind() == reflect.Struct {
		return ferr(errUnsupportedType)
	}
	// create correct parse priority
	var value string
	var exists bool
	priority := priorityOrder()
	for _, p := range priority {
		switch p {
		case FlagPriority:
			value, exists = v.flagV.value()
		case EnvPriority:
			value, exists = v.envV.value()
		case ExternalPriority:
			values := []Value{v}
			owner := v.owner
			for owner != nil && owner.Name() != "" {
				values = append([]Value{owner}, values...)
				owner = owner.parent
			}
			value, exists = v.owner.external.Get(values...)
		case DefaultPriority:
			value, exists = v.defaultV.value()
		}
		if exists {
			traceLogger.Printf("envconf: set variable name=%s value=%s from=%s", v.fullname(), value, p)
			break
		}
	}
	if !exists {
		return ferr(errRequiredFiled)
	}
	// set value
	switch v.tag.Type.Kind() {
	case reflect.Bool:
		i, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.field.SetBool(i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		var err error
		if _, ok := v.field.Interface().(time.Duration); ok {
			var d time.Duration
			d, err = time.ParseDuration(value)
			if err != nil {
				return err
			}
			i = d.Nanoseconds()
		} else {
			i, err = strconv.ParseInt(value, 0, 64)
			if err != nil {
				return err
			}
		}
		v.field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return err
		}
		v.field.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		v.field.SetFloat(i)
	case reflect.String:
		v.field.SetString(value)
	default:
		return ferr(errUnsupportedType)
	}
	return nil
}

func (v *value) String() string {
	return v.Name()
}
