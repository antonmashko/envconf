package envconf

type Priority int

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
