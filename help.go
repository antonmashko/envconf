package envconf

import (
	"fmt"
	"os"
)

// UseCustomHelp override default flag `Usage`
var UseCustomHelp = true

type help struct {
	p *parser
}

func (h *help) usage() {
	fmt.Fprintln(os.Stdout, "Usage:")
	fmt.Fprintln(os.Stdout, "")
	h.print(h.p)
}

func (h *help) print(p *parser) {
	if p == nil {
		return
	}
	for _, v := range p.values {
		h.printValue(v)
	}
	for _, ch := range p.children {
		h.print(ch)
	}
}

func (h *help) printValue(val *value) {
	// TODO: help should be configurable
	defaultValue, _ := val.defaultV.value()
	fmt.Fprintf(os.Stdout, "%s <%s> %s\n", val.fullname(), val.tag.Type, defaultValue)
	fmt.Fprintf(os.Stdout, "\tflag: %s\n", val.flagV.name)
	fmt.Fprintf(os.Stdout, "\tenvironment variable: %s\n", val.envV.name)
	fmt.Fprintf(os.Stdout, "\trequired: %t\n", val.required)
	if val.desc != "" {
		fmt.Fprintf(os.Stdout, "\tdescription: \"%s\"\n", val.desc)
	}
	fmt.Fprintln(os.Stdout)
}
