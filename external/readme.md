# External
This package provides implementations of External option for EnvConf parsing process.

Implement your own external or use any of already implemented

## Usage
```golang
	json := `{"foo":"bar"}`
	tc := struct {
		Foo string `env:"ENV_FOO"`
	}{}
	err := envconf.Parse(&tc, 
        option.WithExternal(jsonconf.Json([]byte(json))))
```
If env variable is it not specified for Foo variable, EnvConf will define it from json input.

## Wrapped external
- json (encoding/json)
- yaml (gopkg.in/yaml.v3)
