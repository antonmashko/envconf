# EnvConf
[![Go Report](https://goreportcard.com/badge/github.com/antonmashko/envconf)](https://goreportcard.com/report/github.com/antonmashko/envconf)
[![GoDoc](http://godoc.org/github.com/antonmashko/envconf?status.svg)](http://godoc.org/github.com/antonmashko/envconf)
[![Build and Test](https://github.com/antonmashko/envconf/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/antonmashko/envconf/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/antonmashko/envconf/branch/master/graph/badge.svg?token=ZdkG2flKKv)](https://codecov.io/gh/antonmashko/envconf)    

EnvConf is a Go package, for parsing configuration values from different sources. 

## Installing
```
go get github.com/antonmashko/envconf
```

## Parse Configs
Usually you need a tag with desire configuration sources and execution of a single function `envconf.Parse` for getting all configuration values into your golang structure.

### Supported Configurations
* command line flags
* environment variables
* default values
* external sources (can be anything that is implementing interface [External](https://pkg.go.dev/github.com/antonmashko/envconf#External))

### Tags
Use tags for getting values from different configuration sources.
- flag - name of flag;   
- env - name of environment variable;
- default - if nothing set this value will be used as field value; 
- required - on `true` checks that configuration exists in `flag` or `env` source;  
- description - field description in help output.
- envconf - only for structs. override struct name for generating configuration name. 

### Supported Types
1. Primitives: `bool`, `string`, all types of `int` and `unit`, `float32`, `float64`, `complex64`, `complex128`;
2. Collections:
	- Array and Slice - comma-separated string can be converted into slice or array. NOTE: if elements in string more than len of array EnvConf will panic with `index out of range`.
	- Map - comma-separated string with a colon-separated key and value can be converted into map. example input: `key1:value1, key2:value2`
3. Golang types:
	- time.Duration;
	- time.Time - using `time.RFC3339` as a time.Parse layout argument;
	- net.IP;
	- url.URL;

### Example
Let's take a look at a simple example. Here we're creating struct with 3 tags for different configuration sources: flag, env, and default value. **NOTE**: It's not necessary to specify tags for each configuration type, add desired only 

```golang
package main

import (
	"fmt"

	"github.com/antonmashko/envconf"
)

type Example struct {
	Field1 string `flag:"flag-name" env:"ENV_VAR_NAME" default:"default-value"`
}

func main() {
	var cfg Example
	if err := envconf.Parse(&cfg); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", cfg)
}
```


**Testing!**
If you want to get set `Field1` from command line flag, use flag name that is set in `flag` tag. 
```bash
$ go run main.go -flag-name="variable-from-flag"
main.Example{Field1:"variable-from-flag"}
```
The same result would be for other configuration types.

### `-help` output
Using envconf will also generate help output with all registered fields and types. Use flag `-help` or `-h` for getting it. 
```bash
$ go run main.go -help

Usage:

Field1 <string> default-value
        flag: flag-name
        environment variable: ENV_VAR_NAME
        required: false
        description: ""
```

## Auto-generating Config Names
EnvConf can generate environment variable name or flag name from golang field path. All you need is to set `*` in specific tag. For environment variables name envconf will use field path in uppercase and underscore as a delimiter. 
Example: 
```golang
type Config struct {
	HTTP struct {
		Addr string `env:"*"`
	}
}
```
Now we can use `HTTP_ADDR` environment variable for defining Addr field. 
The same approach will work for flag. But flag names will be generated in lowercase and the dash will be as a delimiter.

### Overriding Parent struct name for Auto-generation
In case if parent struct name doesn't satisfy for configuration variable name, it can be changed with `envconf` tag.
Example:
```golang
type Config struct {
	HTTP struct {
		Addr string `env:"*"`
	} `envconf:"httpserver"`
}
```
Now we'll get `HTTPSERVER_ADDR` as environment variable name.
See: [EnvConf example](example/main.go)

## External
reading json config
see: [example](example/main.go)

## Options
Options allow intercept into `EnvConf.Parse` process

Name|Option|Description
---|---|---
Read configuration priority|`option.WithPriorityOrder`|Change default parsing priority. Default: *Flag*, *Environment variable*, *External source*, *Default Value*
Log|`option.WithLog`|Enable logging over parsing process. Prints defined and not defined configuration fields
Custom Usage|`option.WithCustomUsage`|Generate usage for `-help` flag from input structure. By default this option is enabled, use `option.WithoutCustomUsage` option
Flag Parsed Callback|`option.WithFlagParsed`|This callback allow to use flags after flag.Parse() and before EnvConf.Define process
