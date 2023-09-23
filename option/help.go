package option

import (
	"flag"
	"fmt"
	"io"
)

type help struct {
	out    io.Writer
	fields []FieldInitializedArg
}

func (h *help) usage() {
	fmt.Fprintln(h.out, "Usage:")
	fmt.Fprintln(h.out, "")
	h.print()
}

func (h *help) print() {
	for _, f := range h.fields {
		h.printValue(f)
	}
}

func (h *help) printValue(f FieldInitializedArg) {
	// TODO: help should be configurable
	fmt.Fprintf(h.out, "%s <%s> %s\n", f.FullName, f.Type.Name(), f.DefaultValue)
	fmt.Fprintf(h.out, "\tflag: %s\n", f.FlagName)
	fmt.Fprintf(h.out, "\tenvironment variable: %s\n", f.EnvName)
	fmt.Fprintf(h.out, "\trequired: %t\n", f.Required)
	if f.Description != "" {
		fmt.Fprintf(h.out, "\tdescription: \"%s\"\n", f.Description)
	}
	fmt.Fprintln(h.out)
}

func (h *help) addField(arg FieldInitializedArg) {
	h.fields = append(h.fields, arg)
}

func (h *help) Apply(opts *Options) {
	opts.usage = h.usage
	opts.onFieldInitialized = h.addField
}

// WithCustomUsage generates usage for -help flag from input struct.
// By default EnvConf uses this function. Use `option.WithoutCustomUsage` to disable it
func WithCustomUsage() ClientOption {
	h := &help{
		out:    flag.CommandLine.Output(),
		fields: make([]FieldInitializedArg, 0),
	}
	return h
}

type disableHelp struct{}

func (disableHelp) Apply(opts *Options) {
	opts.usage = nil
	opts.onFieldInitialized = nil
}

func WithoutCustomUsage() ClientOption {
	return disableHelp{}
}
