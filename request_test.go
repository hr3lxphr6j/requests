package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		writer.WriteHeader(http.StatusOK)
	}))
	_, err := Get(hs.URL)
	assert.NoError(t, err)
}

func TestHead(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodHead, request.Method)
		writer.WriteHeader(http.StatusOK)
	}))
	_, err := Head(hs.URL)
	assert.NoError(t, err)
}

func TestPost(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		writer.WriteHeader(http.StatusOK)
	}))
	_, err := Post(hs.URL)
	assert.NoError(t, err)
}

func TestPut(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPut, request.Method)
		writer.WriteHeader(http.StatusOK)
	}))
	_, err := Put(hs.URL)
	assert.NoError(t, err)
}

func TestPatch(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		writer.WriteHeader(http.StatusOK)
	}))
	_, err := Patch(hs.URL)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		writer.WriteHeader(http.StatusOK)
	}))
	_, err := Delete(hs.URL)
	assert.NoError(t, err)
}

func TestConnect(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodConnect, request.Method)
		writer.WriteHeader(http.StatusOK)
	}))
	_, err := Connect(hs.URL)
	assert.NoError(t, err)
}

func TestOptions(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodOptions, request.Method)
		writer.WriteHeader(http.StatusOK)
	}))
	_, err := Options(hs.URL)
	assert.NoError(t, err)
}

func TestTrace(t *testing.T) {
	hs := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodTrace, request.Method)
		writer.WriteHeader(http.StatusOK)
	}))
	_, err := Trace(hs.URL)
	assert.NoError(t, err)
}
