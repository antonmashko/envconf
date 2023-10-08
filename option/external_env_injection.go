package option

import "strings"

type withExternalEnvInjection func(string) (string, ConfigSource)

func (o withExternalEnvInjection) Apply(opts *Options) {
	opts.externalEnvInjection = o
}

func WithExternalEnvInjection() ClientOption {
	return WithCustomExternalEnvInjection(nil)
}

func WithCustomExternalEnvInjection(f func(string) (string, ConfigSource)) ClientOption {
	if f == nil {
		f = func(value string) (string, ConfigSource) {
			var ok bool
			const prefix = "${"
			value, ok = strings.CutPrefix(value, prefix)
			if !ok {
				return "", NoConfigValue
			}

			const suffix = "}"
			value, ok = strings.CutSuffix(value, suffix)
			if !ok {
				return "", NoConfigValue
			}

			value = strings.TrimSpace(value)
			const envInjectionPrefix = ".env."
			value, ok = strings.CutPrefix(value, envInjectionPrefix)
			if !ok {
				return "", NoConfigValue
			}
			return value, EnvVariable
		}
	}
	return withExternalEnvInjection(f)
}
