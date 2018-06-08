package envconf

type Priority int

func (p Priority) String() string {
	switch p {
	case FlagPriority:
		return "Flag"
	case EnvPriority:
		return "Environment"
	case ConfigFilePriority:
		return "Configuration"
	case DefaultPriority:
		return "Default"
	}
	return ""
}

const (
	FlagPriority Priority = iota
	EnvPriority
	ConfigFilePriority
	DefaultPriority
)

var priorityQueue = map[Priority]int{
	FlagPriority:       0,
	EnvPriority:        1,
	ConfigFilePriority: 2,
	DefaultPriority:    3,
}

func SetPriority(priority ...Priority) {
	for i, p := range priority {
		priorityQueue[p] = i
	}
}

func priorityOrder() []Priority {
	result := make([]Priority, len(priorityQueue))
	for value, index := range priorityQueue {
		result[index] = value
	}
	return result
}
