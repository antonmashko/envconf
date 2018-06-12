package logger

import (
	"log"
	"os"
)

type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Debugln(...interface{})

	Info(...interface{})
	Infof(string, ...interface{})
	Infoln(...interface{})

	Warn(...interface{})
	Warnf(string, ...interface{})
	Warnln(...interface{})

	Error(...interface{})
	Errorf(string, ...interface{})
	Errorln(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
}

type stdoutLogger struct {
	logger *log.Logger
	level  LogLevel
	name   string
}

func New(name string, level LogLevel) Logger {
	return &stdoutLogger{
		logger: log.New(os.Stdout, name, log.Ldate|log.Ltime),
		level:  level,
		name:   name,
	}
}

func (l *stdoutLogger) Debug(args ...interface{}) {
	l.print(Debug, args...)
}

func (l *stdoutLogger) Debugf(format string, args ...interface{}) {
	l.printf(Debug, format, args...)
}

func (l *stdoutLogger) Debugln(args ...interface{}) {
	l.println(Debug, args...)
}

func (l *stdoutLogger) Info(args ...interface{}) {
	l.print(Info, args...)
}

func (l *stdoutLogger) Infof(format string, args ...interface{}) {
	l.printf(Info, format, args...)
}

func (l *stdoutLogger) Infoln(args ...interface{}) {
	l.println(Info, args...)
}

func (l *stdoutLogger) Warn(args ...interface{}) {
	l.print(Warn, args...)
}

func (l *stdoutLogger) Warnf(format string, args ...interface{}) {
	l.printf(Warn, format, args...)
}

func (l *stdoutLogger) Warnln(args ...interface{}) {
	l.println(Warn, args...)
}

func (l *stdoutLogger) Error(args ...interface{}) {
	l.print(Error, args...)
}

func (l *stdoutLogger) Errorf(format string, args ...interface{}) {
	l.printf(Error, format, args...)
}

func (l *stdoutLogger) Errorln(args ...interface{}) {
	l.println(Error, args...)
}

func (l *stdoutLogger) Fatal(args ...interface{}) {
	l.print(Fatal, args...)
}

func (l *stdoutLogger) Fatalf(format string, args ...interface{}) {
	l.printf(Fatal, format, args...)
}

func (l *stdoutLogger) Fatalln(args ...interface{}) {
	l.println(Fatal, args...)
}

func (l *stdoutLogger) print(level LogLevel, args ...interface{}) {
	if level > l.level {
		return
	}
	l.logger.Print(append([]interface{}{level, l.logger.Prefix()}, args...)...)
}

func (l *stdoutLogger) printf(level LogLevel, format string, args ...interface{}) {
	if level > l.level {
		return
	}
	l.logger.Printf(format, append([]interface{}{level, l.logger.Prefix()}, args...)...)
}

func (l *stdoutLogger) println(level LogLevel, args ...interface{}) {
	if level > l.level {
		return
	}
	l.logger.Println(append([]interface{}{level, l.logger.Prefix()}, args...)...)
}
