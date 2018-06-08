package envconf

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"strconv"
	"strings"
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
	name  string
	value string
}

func (f *flagV) define(tag reflect.StructField, usage string) {
	f.name = tag.Tag.Get(tagFlag)
	if f.name == tagNotDefined {
		f.name = tagIgnored
	} else if strings.ToLower(f.name) == valDefault {
		f.name = strings.ToLower(tag.Name)
	}
	if f.name != tagIgnored {
		flag.StringVar(&f.value, f.name, "", usage)
	}
}

type envV struct {
	name  string
	value string
}

func (e *envV) define(tag reflect.StructField) {
	//define env name
	e.name = tag.Tag.Get(tagEnv)
	if e.name == tagNotDefined {
		e.name = tagIgnored
	} else if e.name == valDefault {
		e.name = strings.ToUpper(tag.Name)
	}
	//get value
	if e.name != tagIgnored {
		os.Getenv(e.name)
	}
}

type value struct {
	owner        *parser
	field        reflect.Value
	tag          reflect.StructField
	flagV        flagV  // flag value
	envV         envV   // env value
	required     bool   // if it defined true, value should be defined
	defaultValue string // default value
	desc         string // description
}

func newValue(field reflect.Value, tag reflect.StructField) *value {
	v := &value{field: field, tag: tag}
	// Parse description
	v.desc = tag.Tag.Get(tagDescription)
	(&v.flagV).define(tag, v.desc)
	(&v.envV).define(tag)
	// Parse required
	rq := tag.Tag.Get(tagRequired)
	v.required, _ = strconv.ParseBool(rq)
	// Parse default value
	v.defaultValue = tag.Tag.Get(tagDefault)
	return v
}

func (v *value) name() string {
	return v.owner.Path() + Separator + v.tag.Name
}

func (v *value) find() error {
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
	var value string
	// create correct parse priority
	for p, f := range priorityQueue {
		switch p {
		case FlagPriority:
			value = v.flagV.value
		case EnvPriority:
			value = v.envV.value
		case ConfigFilePriority:
			break // TODO: handle it correct
		case DefaultPriority:
			value = v.defaultValue
		}
		if value != "" {
			break
		}
	}
	if value == "" {
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
		i, err := strconv.ParseInt(value, 0, 64)
		if err != nil {
			return err
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
