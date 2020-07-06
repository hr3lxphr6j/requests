package requests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
}

func NewResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

func (r *Response) StdResponse() *http.Response {
	return r.Response
}

func (r *Response) Bytes() ([]byte, error) {
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

func (r *Response) Text() (string, error) {
	b, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *Response) JSON(i interface{}) error {
	return json.NewDecoder(r.Body).Decode(i)
}
