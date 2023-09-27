package option

import "github.com/antonmashko/envconf/external"

type extOpt struct {
	ext external.External
}

func (o extOpt) Apply(opts *Options) {
	opts.external = o.ext
}

func WithExternal(e external.External) ClientOption {
	return extOpt{ext: e}
}
