package option

import "strings"

type withExternalEnvInjection func(string) (string, ConfigSource)

func (f withExternalEnvInjection) Apply(opts *Options) {
	opts.externalInjection = f
}

// WithExternalInjection allows to inject variables from an external source
// use format ${ .env.<ENV_VAR_NAME> } to make an env variable lookup while getting a field from the external source
func WithExternalInjection() ClientOption {
	return WithCustomExternalInjection(nil)
}

func WithCustomExternalInjection(f func(string) (string, ConfigSource)) ClientOption {
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
