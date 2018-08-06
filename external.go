package envconf

type External interface {
	Get(...Value) (string, bool)
}

type emptyExt struct{}

func (c *emptyExt) Get(v ...Value) (string, bool) { return "", false }
