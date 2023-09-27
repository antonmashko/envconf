package external

// External config source
// Implementation of this interface should be able to `Unmarshal` data into map[string]interface{},
// where interface{} should be also same map type for the nested structures
type External interface {
	// TagName is key name in golang struct tag (json, yaml, toml etc.).
	TagName() []string
	// Unmarshal parses the external data and stores the result
	// in the value pointed to by v.
	// Usually, it just wraps the existing `Unmarshal` function of third-party libraries
	Unmarshal(v interface{}) error
}
