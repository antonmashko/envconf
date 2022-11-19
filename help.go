package envconf

import (
	"fmt"
	"os"
)

// UseCustomHelp override default flag `Usage`
var UseCustomHelp = true

type help struct {
	fields []*primitiveType
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

func (h *help) printValue(f *primitiveType) {
	// TODO: help should be configurable
	defaultValue, _ := f.def.Value()
	fmt.Fprintf(os.Stdout, "%s <%s> %s\n", fullname(f), f.tag.Type.Name(), defaultValue)
	fmt.Fprintf(os.Stdout, "\tflag: %s\n", f.flag.Name())
	fmt.Fprintf(os.Stdout, "\tenvironment variable: %s\n", f.env.Name())
	fmt.Fprintf(os.Stdout, "\trequired: %t\n", f.required)
	if f.desc != "" {
		fmt.Fprintf(os.Stdout, "\tdescription: \"%s\"\n", f.desc)
	}
	fmt.Fprintln(os.Stdout)
}
