package envconf

type Priority int

func (p Priority) String() string {
	switch p {
	case FlagPriority:
		return "Flag"
	case EnvPriority:
		return "Environment"
	case ExternalPriority:
		return "External"
	case DefaultPriority:
		return "Default"
	}
	return ""
}

const (
	FlagPriority Priority = iota
	EnvPriority
	ExternalPriority
	DefaultPriority
)

var priorityQueue = map[Priority]int{
	FlagPriority:     0,
	EnvPriority:      1,
	ExternalPriority: 2,
	DefaultPriority:  3,
}

// SetPriority can override default priority queue.
// Default priority queue is: Flag, Environment variable, External source, Default value.
func SetPriority(priority ...Priority) {
	if len(priority) == 0 {
		return
	}
	priorityQueue = make(map[Priority]int)
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
