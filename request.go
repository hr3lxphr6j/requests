package requests

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

func NewRequest(method, url string, opts ...RequestOption) (*http.Request, error) {
	options := NewOptions()
	for _, opt := range opts {
		opt(options)
		if options.Err != nil {
			return nil, options.Err
		}
	}
	if options.multipartWriter != nil {
		if err := options.multipartWriter.Close(); err != nil {
			return nil, fmt.Errorf("write multipart error %w", err)
		}
	}

	ctx := context.Background()
	if !options.Deadline.IsZero() {
		ctx, _ = context.WithDeadline(ctx, options.Deadline)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, options.Body)
	if err != nil {
		return nil, err
	}

	// set headers
	for k, v := range options.Headers {
		switch val := v.(type) {
		case string:
			req.Header.Set(k, val)
		case fmt.Stringer:
			req.Header.Set(k, val.String())
		case nil:
			req.Header.Del(k)
		default:
			return nil, fmt.Errorf("value of header [%s] must be string or nil, but %s", k, reflect.TypeOf(v))
		}
	}

	// set queries
	req.URL.RawQuery = options.Queries.Encode()

	// set cookies
	for k, v := range options.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	return req, nil
}

func Get(url string, opts ...RequestOption) (*Response, error) {
	return DefaultSession.Get(url, opts...)
}

func Head(url string, opts ...RequestOption) (*Response, error) {
	return DefaultSession.Get(url, opts...)
}

func Post(url string, opts ...RequestOption) (*Response, error) {
	return DefaultSession.Get(url, opts...)
}

func Put(url string, opts ...RequestOption) (*Response, error) {
	return DefaultSession.Get(url, opts...)
}

func Patch(url string, opts ...RequestOption) (*Response, error) {
	return DefaultSession.Get(url, opts...)
}

func Delete(url string, opts ...RequestOption) (*Response, error) {
	return DefaultSession.Get(url, opts...)
}

func Connect(url string, opts ...RequestOption) (*Response, error) {
	return DefaultSession.Get(url, opts...)
}

func Options(url string, opts ...RequestOption) (*Response, error) {
	return DefaultSession.Get(url, opts...)
}

func Trace(url string, opts ...RequestOption) (*Response, error) {
	return DefaultSession.Get(url, opts...)
}
