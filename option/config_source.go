package option

type ConfigSource int

const (
	NoConfigValue ConfigSource = -1
	FlagVariable  ConfigSource = 1 << (iota - 1)
	EnvVariable
	ExternalSource
	DefaultValue
)

func (s ConfigSource) String() string {
	switch s {
	case FlagVariable:
		return "Flag"
	case EnvVariable:
		return "Environment"
	case ExternalSource:
		return "External"
	case DefaultValue:
		return "Default"
	}
	return ""
}

type priorityOrder []ConfigSource

func (p priorityOrder) Apply(opts *Options) {
	opts.priorityOrder = p
}

// WithPriorityOrder overrides default priority order, with an order from function argument.
// Default priority order is: Flag, Environment variable, External source, Default value.
func WithPriorityOrder(s ...ConfigSource) ClientOption {
	defaultOrder := []ConfigSource{
		FlagVariable, EnvVariable, ExternalSource, DefaultValue,
	}
	if len(s) == 0 {
		return priorityOrder(defaultOrder)
	}
	po := make(map[ConfigSource]int)
	var idx int
	for _, p := range s {
		if p != FlagVariable && p != EnvVariable &&
			p != ExternalSource && p != DefaultValue {
			continue
		}
		if _, ok := po[p]; !ok {
			po[p] = idx
			idx++
		}
	}
	if len(po) == 0 {
		return priorityOrder(defaultOrder)
	}
	result := make([]ConfigSource, len(po))
	for s, idx := range po {
		result[idx] = s
	}
	return priorityOrder(result)
}
