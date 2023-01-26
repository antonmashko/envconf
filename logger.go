package envconf

import (
	"io/ioutil"
	"log"
)

var debugLogger Logger = &logger{l: log.New(ioutil.Discard, "", log.Ltime)}

type Logger interface {
	Printf(string, ...interface{})
	Println(...interface{})
}

type logger struct {
	l *log.Logger
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.l.Printf(format, args...)
}

func (l *logger) Println(args ...interface{}) {
	l.l.Println(args...)
}
