package envconf

import (
	"fmt"
	"io"
)

type Logger interface {
	Printf(string, ...interface{})
	Println(...interface{})
}

type logger struct {
	w io.Writer
}

func (l *logger) Printf(format string, args ...interface{}) {
	fmt.Fprintf(l.w, format, args...)
}

func (l *logger) Println(args ...interface{}) {
	fmt.Fprintln(l.w, args...)
}
