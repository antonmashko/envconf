package option

import (
	"io"
	"log"
	"regexp"
)

type Printer interface {
	Print(...interface{})
}

type Logger struct {
	Printer
	HideSecrets      bool
	SecretMatchRegex string
}

func (l *Logger) printDefined(arg FieldDefinedArg) {
	var v interface{} = arg.Value
	if l.HideSecrets && regexp.MustCompile(l.SecretMatchRegex).MatchString(arg.Name) {
		v = "******"
	}
	l.Print("field=\"", arg.FullName, "\" value=\"", v,
		"\" source=\"", arg.Source.String(), "\"")
}

func (l *Logger) printErr(arg FieldDefineErrorArg) {
	l.Print("field=\"", arg.FullName, "\" err=\"", arg.Err, "\"")
}

func (l *Logger) Apply(opts *Options) {
	opts.OnFieldDefined = append(opts.OnFieldDefined, l.printDefined)
	opts.OnFieldDefineErr = append(opts.OnFieldDefineErr, l.printErr)
}

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
