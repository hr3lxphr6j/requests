package requests

import (
	"io"
)

type Options struct {
	Headers map[string]string
	Queries map[string]string
	Cookies map[string]string
	Body    io.Reader

	err error
}

type Option func(o *Options)

func Header(k, v string) Option {
	return func(o *Options) {
		o.Headers[k] = v
	}
}

func Headers(m map[string]string) Option {
	return func(o *Options) {
		for k, v := range m {
			o.Headers[k] = v
		}
	}
}

func Query(k, v string) Option {
	return func(o *Options) {
		o.Queries[k] = v
	}
}

func () {

}
