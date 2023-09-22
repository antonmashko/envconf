package option

import (
	"fmt"
	"os"
)

type help struct {
	fields []FieldInitializedArg
}

func (h *help) usage() {
	fmt.Fprintln(os.Stdout, "Usage:")
	fmt.Fprintln(os.Stdout, "")
	h.print()
}

func (h *help) print() {
	for _, f := range h.fields {
		h.printValue(f)
	}
}

func (h *help) printValue(f FieldInitializedArg) {
	// TODO: help should be configurable
	fmt.Fprintf(os.Stdout, "%s <%s> %s\n", f.FullName, f.Type.Name(), f.DefaultValue)
	fmt.Fprintf(os.Stdout, "\tflag: %s\n", f.FlagName)
	fmt.Fprintf(os.Stdout, "\tenvironment variable: %s\n", f.EnvName)
	fmt.Fprintf(os.Stdout, "\trequired: %t\n", f.Required)
	if f.Description != "" {
		fmt.Fprintf(os.Stdout, "\tdescription: \"%s\"\n", f.Description)
	}
	fmt.Fprintln(os.Stdout)
}

func (h *help) addField(arg FieldInitializedArg) {
	h.fields = append(h.fields, arg)
}

func (h *help) Apply(opts *Options) {
	opts.Usage = h.usage
	opts.OnFieldInitialized = append(opts.OnFieldInitialized, h.addField)
}

func WithFlagUsage() ClientOption {
	h := &help{
		fields: make([]FieldInitializedArg, 0),
	}
	return h
}
