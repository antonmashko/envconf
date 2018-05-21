# EnvConf
## What is EnvConf?  
EnvConf is a Go package, for parsing configuration values from different sources. 

## How it works?
Create golang struct with following tags
``` golang
type Example struct {
    Field1 string `flag:"field" env:"EXAMPLE_FIELD" "default:"example_value" required:"false" description:"this is exaple configurable field"`
}
```