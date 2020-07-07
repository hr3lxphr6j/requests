package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse_StdResponse(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	resp, err := Get(hs.URL)
	assert.NoError(t, err)
	assert.NotNil(t, resp.StdResponse())
}

func TestResponse_Text(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("foo"))
	}))
	resp, err := Get(hs.URL)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	s, err := resp.Text()
	assert.NoError(t, err)
	assert.Equal(t, "foo", s)
}

func TestResponse_JSON(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(`{"a":"b", "c": "d"}`))
	}))
	resp, err := Get(hs.URL)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	m := map[string]interface{}{}
	assert.NoError(t, resp.JSON(&m))
	assert.Equal(t, map[string]interface{}{"a": "b", "c": "d"}, m)
}
