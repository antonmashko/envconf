package envconf

// External config source
type External interface {
	// Get string value from values chain(from parent to child)
	Get(...Value) (interface{}, bool)
	//
	Unmarshal(interface{}) error
}

type emptyExt struct{}

func (c *emptyExt) Get(v ...Value) (interface{}, bool) { return nil, false }

func (c *emptyExt) Unmarshal(v interface{}) error { return nil }
