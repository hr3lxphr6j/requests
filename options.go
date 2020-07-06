package requests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	HeaderContentType   = "Content-Type"
	HeaderUserAgent     = "User-Agent"
	HeaderReferer       = "Referer"
	HeaderAuthorization = "Authorization"

	ContentTypeJSON = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"
)

type RequestOptions struct {
	Headers  map[string]interface{}
	Queries  url.Values
	Cookies  map[string]string
	Body     io.Reader
	Deadline time.Time

	multipartWriter *multipart.Writer

	Err error
}

func (o *RequestOptions) getMultipartWriter() *multipart.Writer {
	if o.multipartWriter == nil {
		Body(new(bytes.Buffer))(o)
		o.multipartWriter = multipart.NewWriter(o.Body.(*bytes.Buffer))
	}
	return o.multipartWriter
}

func NewOptions() *RequestOptions {
	return &RequestOptions{
		Headers: map[string]interface{}{},
		Queries: url.Values{},
		Cookies: map[string]string{},
	}
}

type RequestOption func(o *RequestOptions)

func Deadline(t time.Time) RequestOption {
	return func(o *RequestOptions) {
		o.Deadline = t
	}
}

func Timeout(d time.Duration) RequestOption {
	return func(o *RequestOptions) {
		o.Deadline = time.Now().Add(d)
	}
}

func Header(k string, v interface{}) RequestOption {
	return func(o *RequestOptions) {
		o.Headers[k] = v
	}
}

func Headers(m map[string]interface{}, replace ...bool) RequestOption {
	return func(o *RequestOptions) {
		if len(replace) == 1 && replace[1] {
			o.Headers = m
			return
		}
		for k, v := range m {
			o.Headers[k] = v
		}
	}
}

func UserAgent(ua string) RequestOption {
	return func(o *RequestOptions) {
		Header(HeaderUserAgent, ua)(o)
	}
}

func ContentType(ct string) RequestOption {
	return func(o *RequestOptions) {
		Header(HeaderContentType, ct)(o)
	}
}

func Referer(r string) RequestOption {
	return func(o *RequestOptions) {
		Header(HeaderReferer, r)(o)
	}
}

func Authorization(a string) RequestOption {
	return func(o *RequestOptions) {
		Header(HeaderAuthorization, a)(o)
	}
}

func BasicAuth(username, password string) RequestOption {
	return func(o *RequestOptions) {
		Authorization(base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", username, password))))
	}
}

func Cookie(k, v string) RequestOption {
	return func(o *RequestOptions) {
		o.Cookies[k] = v
	}
}

func Cookies(m map[string]string, replace ...bool) RequestOption {
	return func(o *RequestOptions) {
		if len(replace) == 1 && replace[1] {
			o.Cookies = m
			return
		}
		for k, v := range m {
			o.Cookies[k] = v
		}
	}
}

func Query(k, v string) RequestOption {
	return func(o *RequestOptions) {
		o.Queries.Set(k, v)
	}
}

func Queries(m map[string]string, replace ...bool) RequestOption {
	return func(o *RequestOptions) {
		if len(replace) == 1 && replace[1] {
			o.Queries = url.Values{}
		}
		for k, v := range m {
			o.Queries.Set(k, v)
		}
	}
}

func QueriesFromValue(v url.Values) RequestOption {
	return func(o *RequestOptions) {
		o.Queries = v
	}
}

func Body(r io.Reader) RequestOption {
	return func(o *RequestOptions) {
		o.Body = r
	}
}

func JSON(i interface{}) RequestOption {
	return func(o *RequestOptions) {
		ContentType(ContentTypeJSON)(o)
		Body(new(bytes.Buffer))(o)
		if err := json.NewEncoder(o.Body.(*bytes.Buffer)).Encode(i); err != nil {
			o.Err = fmt.Errorf("json option error %w", err)
		}
	}
}

func Form(m map[string]string) RequestOption {
	return func(o *RequestOptions) {
		ContentType(ContentTypeForm)(o)
		value := url.Values{}
		for k, v := range m {
			value.Set(k, v)
		}
		Body(strings.NewReader(value.Encode()))(o)
	}
}

func File(fieldName string, file *os.File) RequestOption {
	return func(o *RequestOptions) {
		mw := o.getMultipartWriter()
		w, err := mw.CreateFormFile(fieldName, filepath.Base(file.Name()))
		if err != nil {
			o.Err = fmt.Errorf("file option error %w", err)
			return
		}
		_, err = io.Copy(w, file)
		if err != nil {
			o.Err = fmt.Errorf("file option error %w", err)
			return
		}
		ContentType(mw.FormDataContentType())(o)
	}
}

func MultipartFieldString(fieldName, value string) RequestOption {
	return func(o *RequestOptions) {
		mw := o.getMultipartWriter()
		if err := mw.WriteField(fieldName, value); err != nil {
			o.Err = fmt.Errorf("multipartFieldString option error %w", err)
			return
		}
		ContentType(mw.FormDataContentType())(o)
	}
}

func MultipartField(fieldName string, r io.Reader) RequestOption {
	return func(o *RequestOptions) {
		mw := o.getMultipartWriter()
		w, err := mw.CreateFormField(fieldName)
		if err != nil {
			o.Err = fmt.Errorf("multipartField option error %w", err)
			return
		}
		_, err = io.Copy(w, r)
		if err != nil {
			o.Err = fmt.Errorf("file option error %w", err)
			return
		}
		ContentType(mw.FormDataContentType())(o)
	}
}
