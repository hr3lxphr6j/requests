package requests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeadline(t *testing.T) {
	done := make(chan struct{})

	hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-done
		w.WriteHeader(http.StatusOK)
	}))

	go func() {
		defer close(done)
		_, err := Request(http.MethodPost, hs.URL, Deadline(time.Now().Add(time.Second)))
		assert.Error(t, err)
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("timeout")
	case <-done:
	}
}

func TestTimeout(t *testing.T) {
	done := make(chan struct{})

	hs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-done
		w.WriteHeader(http.StatusOK)
	}))

	go func() {
		defer close(done)
		_, err := Request(http.MethodPost, hs.URL, Timeout(time.Second))
		assert.Error(t, err)
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("timeout")
	case <-done:
	}
}

func TestHeader(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://localhost",
		Header("key_1", "value_1"),
		Header("key_2", "value_2"),
		Header("key_2", nil),
		Header("key_1", "value_3"),
	)
	assert.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"Key_1": {"value_3"},
	}, map[string][]string(req.Header))
}

func TestHeaderWithBadValue(t *testing.T) {
	_, err := NewRequest(http.MethodGet, "http://localhost",
		Header("key_1", func() {}),
	)
	assert.Error(t, err)
}

func TestHeaders(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://localhost",
		Header("foo", "bar"),
		Headers(map[string]interface{}{
			"key_1": "value_1",
			"key_2": "value_2",
		}, true),
		Headers(map[string]interface{}{
			"key_1": "value_3",
			"key_2": nil,
		}),
	)
	assert.NoError(t, err)
	assert.Equal(t, map[string][]string{
		"Key_1": {"value_3"},
	}, map[string][]string(req.Header))
}

func TestHeadersWithBadvalue(t *testing.T) {
	_, err := NewRequest(http.MethodGet, "http://localhost",
		Headers(map[string]interface{}{
			"key_1": "foo",
			"key_2": func() {},
		}, true),
	)
	assert.Error(t, err)
}

func TestUserAgent(t *testing.T) {
	ua := "requests-0.0.1"
	req, err := NewRequest(http.MethodGet, "http://localhost", UserAgent(ua))
	assert.NoError(t, err)
	assert.Equal(t, ua, req.Header.Get(HeaderUserAgent))
}

func TestContentType(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://localhost", ContentType(ContentTypeJSON))
	assert.NoError(t, err)
	assert.Equal(t, ContentTypeJSON, req.Header.Get(HeaderContentType))
}

func TestReferer(t *testing.T) {
	referer := "http://localhost"
	req, err := NewRequest(http.MethodGet, "http://localhost", Referer(referer))
	assert.NoError(t, err)
	assert.Equal(t, referer, req.Header.Get(HeaderReferer))
}

func TestAuthorization(t *testing.T) {
	auth := "foo"
	req, err := NewRequest(http.MethodGet, "http://localhost", Authorization(auth))
	assert.NoError(t, err)
	assert.Equal(t, auth, req.Header.Get(HeaderAuthorization))
}

func TestBasicAuth(t *testing.T) {
	username := "username"
	password := "123456"
	req, err := NewRequest(http.MethodGet, "http://localhost", BasicAuth(username, password))
	assert.NoError(t, err)
	u, p, ok := req.BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, username, u)
	assert.Equal(t, password, p)
}

func TestCookie(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://localhost", Cookie("foo", "bar"))
	assert.NoError(t, err)
	c, err := req.Cookie("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", c.Value)
}

func TestCookies(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://localhost",
		Cookies(map[string]string{
			"foo": "bar",
		}),
		Cookies(map[string]string{
			"a": "b",
		}, true),
	)
	assert.NoError(t, err)
	assert.EqualValues(t, []*http.Cookie{
		{Name: "a", Value: "b"},
	}, req.Cookies())
}

func TestQuery(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://localhost?a=b&c=d", Query("e", "f"))
	assert.NoError(t, err)
	assert.EqualValues(t, map[string][]string{
		"a": {"b"},
		"c": {"d"},
		"e": {"f"},
	}, req.URL.Query())
}

func TestQueries(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://localhost?a=b&c=d",
		Queries(map[string]string{"foo": "bar"}),
		Queries(map[string]string{"e": "f"}, true),
	)
	assert.NoError(t, err)
	assert.EqualValues(t, map[string][]string{
		"a": {"b"},
		"c": {"d"},
		"e": {"f"},
	}, req.URL.Query())
}

func TestQueriesFromValue(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://localhost?a=b&c=d",
		QueriesFromValue(map[string][]string{
			"e": {"f"},
		}),
	)
	assert.NoError(t, err)
	assert.EqualValues(t, map[string][]string{
		"a": {"b"},
		"c": {"d"},
		"e": {"f"},
	}, req.URL.Query())
}

func TestBody(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://localhost?a=b&c=d", Body(strings.NewReader("foo")))
	assert.NoError(t, err)
	b, err := ioutil.ReadAll(req.Body)
	assert.NoError(t, err)
	assert.Equal(t, "foo", string(b))
}

func TestJSON(t *testing.T) {
	m := map[string]interface{}{
		"a": "b",
		"c": []interface{}{"d", "e"},
	}
	req, err := NewRequest(http.MethodPost, "http://localhost", JSON(m))
	assert.NoError(t, err)
	assert.Equal(t, ContentTypeJSON, req.Header.Get(HeaderContentType))
	b, err := ioutil.ReadAll(req.Body)
	assert.NoError(t, err)
	_m := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal(b, &_m))
	assert.EqualValues(t, m, _m)
}

func TestJSONWithBadValue(t *testing.T) {
	_, err := NewRequest(http.MethodPost, "http://localhost", JSON(make(chan struct{})))
	assert.Error(t, err)
}

func TestForm(t *testing.T) {
	req, err := NewRequest(http.MethodPost, "http://localhost", Form(map[string]string{
		"a": "b",
		"c": "d",
	}))
	assert.NoError(t, err)
	assert.Equal(t, ContentTypeForm, req.Header.Get(HeaderContentType))
	b, err := ioutil.ReadAll(req.Body)
	assert.NoError(t, err)
	assert.Equal(t, "a=b&c=d", string(b))
}

func TestFile(t *testing.T) {
	fileContent := []byte("foo\n")
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.True(t, strings.HasPrefix(request.Header.Get(HeaderContentType), "multipart/form-data"))
		f, h, err := request.FormFile("file")
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(h.Filename, "test"))
		b, err := ioutil.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(t, fileContent, b)
	}))
	f, err := ioutil.TempFile("", "test")
	assert.NoError(t, err)
	defer os.Remove(f.Name()) // clean up
	_, err = f.Write(fileContent)
	assert.NoError(t, err)
	assert.NoError(t, f.Sync())
	_, err = f.Seek(0, 0)
	assert.NoError(t, err)

	_, err = Post(hs.URL, File("file", f))
	assert.NoError(t, err)
}

func TestMultipartFieldString(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.True(t, strings.HasPrefix(request.Header.Get(HeaderContentType), "multipart/form-data"))
		assert.Equal(t, "bar", request.FormValue("foo"))
	}))
	_, err := Post(hs.URL, MultipartFieldString("foo", "bar"))
	assert.NoError(t, err)
}

func TestMultipartFieldBytes(t *testing.T) {
	content := []byte("bar")
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.True(t, strings.HasPrefix(request.Header.Get(HeaderContentType), "multipart/form-data"))
		assert.EqualValues(t, content, request.FormValue("foo"))
	}))
	_, err := Post(hs.URL, MultipartFieldBytes("foo", content))
	assert.NoError(t, err)
}
