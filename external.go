package envconf

// External config source
type External interface {
	// Get string value from values chain(from parent to child)
	Get(...Value) (string, bool)
}

type emptyExt struct{}

func (c *emptyExt) Get(v ...Value) (string, bool) { return "", false }
