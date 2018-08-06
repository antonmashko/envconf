# EnvConf
[![Go Report](https://goreportcard.com/badge/github.com/antonmashko/envconf)](https://goreportcard.com/report/github.com/antonmashko/envconf)
[![GoDoc](http://godoc.org/github.com/antonmashko/envconf?status.svg)](http://godoc.org/github.com/antonmashko/envconf)
[![Build Status](https://travis-ci.org/antonmashko/envconf.svg)](https://travis-ci.org/antonmashko/envconf)
[![Codecov](https://img.shields.io/codecov/c/github/antonmashko/envconf.svg)](https://codecov.io/gh/antonmashko/envconf)
    
Installing
```
go get github.com/antonmashko/envconf
```

## What is EnvConf?  
EnvConf is a Go package, for parsing configuration values from different sources. 

## How it works?
- Import envconf package:
``` golang
import "github.com/antonmashko/envconf"
```

- Create golang struct with following tags (NOTE: each of fields can be omitted, here used all just for example)
``` golang
type Example struct {
    Field1 string `flag:"field" env:"EXAMPLE_FIELD" default:"example_value" required:"false" description:"this is exaple configurable field"`
}
```

- Parse outside values to `Exaple` struct
``` golang
func main() {
    var e Example
    envconf.Parse(&e)
    println(e.Field1)
}
```

- Testings application with different passing params:

a) flag 
> ./example -field='flag example'   
> flag example  

b) environment variable 
> EXAMPLE_FIELD='env variable example' ./example    
> env variable example  

c) default value    
> ./example     
> example_value 

## How to get all registered fields?
To print all values which can be parsed use flag `-help` or `-h`:
> ./example -help   
```
Usage:

Field1 <string> example_value
        flag: field
        environment variable: EXAMPLE_FIELD
        required: false
        description: "this is exaple configurable field"
```

## Tags description:
**Priority**:   
```
1) Flag 
2) Environment variable 
3) External source
4) Default value
```

**Tags**: 
- flag - name of flag for field [must be unique];   
- env - name of environment variable for field [must be unique];
- default - if nothing set this value will be used; 
- required - on `true` validate value from `flag` or `env` source;  
- description - description in `help`   
