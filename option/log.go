package option

import (
	"io"
	"log"
	"regexp"
	"strings"
)

type Printer interface {
	Print(...interface{})
}

type Logger struct {
	Printer
	// Enable secrets hide
	HideSecrets bool
	// Regex for match secrets by field name in lower case
	SecretMatchRegex string
}

func (l *Logger) printDefined(arg FieldDefinedArg) {
	var v interface{} = arg.Value
	if l.HideSecrets && regexp.MustCompile(l.SecretMatchRegex).MatchString(strings.ToLower(arg.Name)) {
		v = "******"
	}
	l.Print("field=\"", arg.FullName, "\" value=\"", v, "\" type=\"", arg.Type.String(),
		"\" source=\"", arg.Source.String(), "\"")
}

func (l *Logger) printErr(arg FieldDefineErrorArg) {
	l.Print("field=\"", arg.FullName, "\" err=\"", arg.Err, "\"")
}

func (l *Logger) Apply(opts *Options) {
	opts.onFieldDefined = l.printDefined
	opts.onFieldDefineErr = l.printErr
}

// SetLogger define printer.
// This logger will print defined values and errors to Printer.Print method.
// By Default all secrets will be hidden
func WithLog(p Printer) *Logger {
	if p == nil {
		p = log.New(io.Discard, "envconf", log.Ldate|log.Ltime)
	}
	const secretRegex = "\\b(?:password|token|secret|key|auth|passphrase|private[_ ]?key|api[_ ]?key|credit[_ ]?card)\\b"
	return &Logger{
		Printer:          p,
		HideSecrets:      true,
		SecretMatchRegex: secretRegex,
	}
}
