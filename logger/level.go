package logger

type LogLevel int

const (
	Fatal LogLevel = iota
	Error
	Warn
	Info
	Debug
)

func (l LogLevel) String() string {
	switch l {
	case Fatal:
		return "FATAL"
	case Error:
		return "ERROR"
	case Warn:
		return "WARN"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	}
	return ""
}
