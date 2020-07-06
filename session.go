package requests

import (
	"net/http"
)

type Session struct {
	*http.Client
}

var DefaultSession = NewSession(http.DefaultClient)

func NewSession(c *http.Client) *Session {
	return &Session{Client: c}
}

func (s *Session) Do(r *http.Request) (*Response, error) {
	resp, err := s.Client.Do(r)
	if err != nil {
		return nil, err
	}
	return NewResponse(resp), nil
}

func (s *Session) DoRequest(method, url string, opts ...RequestOption) (*Response, error) {
	req, err := NewRequest(method, url, opts...)
	if err != nil {
		return nil, err
	}
	return s.Do(req)
}

func (s *Session) Get(url string, opts ...RequestOption) (*Response, error) {
	return s.DoRequest(http.MethodGet, url, opts...)
}

func (s *Session) Head(url string, opts ...RequestOption) (*Response, error) {
	return s.DoRequest(http.MethodHead, url, opts...)
}

func (s *Session) Post(url string, opts ...RequestOption) (*Response, error) {
	return s.DoRequest(http.MethodPost, url, opts...)
}

func (s *Session) Put(url string, opts ...RequestOption) (*Response, error) {
	return s.DoRequest(http.MethodPut, url, opts...)
}

func (s *Session) Patch(url string, opts ...RequestOption) (*Response, error) {
	return s.DoRequest(http.MethodPatch, url, opts...)
}

func (s *Session) Delete(url string, opts ...RequestOption) (*Response, error) {
	return s.DoRequest(http.MethodDelete, url, opts...)
}

func (s *Session) Connect(url string, opts ...RequestOption) (*Response, error) {
	return s.DoRequest(http.MethodConnect, url, opts...)
}

func (s *Session) Options(url string, opts ...RequestOption) (*Response, error) {
	return s.DoRequest(http.MethodOptions, url, opts...)
}

func (s *Session) Trace(url string, opts ...RequestOption) (*Response, error) {
	return s.DoRequest(http.MethodTrace, url, opts...)
}
