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
	for _, ch := range p.childs {
		h.print(ch)
	}
}

func (h *help) printValue(val *value) {
	//TODO: ащкьфе shound be confirugable
	fmt.Fprintf(os.Stdout, "%s <%s> %s\n", val.name(), val.tag.Type, val.def)
	fmt.Fprintf(os.Stdout, "\tflag: %s\n", val.flagv)
	fmt.Fprintf(os.Stdout, "\tenvironment variable: %s\n", val.envv)
	fmt.Fprintf(os.Stdout, "\trequired: %t\n", val.required)
	fmt.Fprintf(os.Stdout, "\tdescription: \"%s\"\n\n", val.desc)
}
