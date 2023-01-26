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

### Supported configuration sources
* command line flags
* environment variables
* default values
* external sources (can be anything that is implementing interface [External](https://pkg.go.dev/github.com/antonmashko/envconf#External))

### Supported tags
**Tags**: 
- flag - name of flag;   
- env - name of environment variable;
- default - if nothing set this value will be used as field value; 
- required - on `true` checks that configuration exists in `flag` or `env` source;  
- description - field description in help output.


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

## Configuration Priority
**Priority**:   
```
1) Flag 
2) Environment variable 
3) External source
4) Default value
```
