package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hailkomputer/kvicksand/internal/api"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type cacheStoreMock struct {
	setFn func(key, value string)
	getFn func(key string) (string, bool)
}

func (m *cacheStoreMock) Set(key, value string) {
	if m != nil && m.setFn != nil {
		m.setFn(key, value)
	}
}

func (m *cacheStoreMock) Get(key string) (string, bool) {
	if m != nil && m.getFn != nil {
		return m.getFn(key)
	}
	return "value", true
}

func TestIndex(t *testing.T) {
	cache := &cacheStoreMock{}
	apiHandler := api.NewApiHandler(cache)
	tests := []struct {
		name       string
		cacheStore *cacheStoreMock
		req        *http.Request
		code       int
	}{
		{
			name:       "should return 200",
			cacheStore: nil,
			req:        &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/"}, Body: nil},
			code:       http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			apiHandler.Router.ServeHTTP(recorder, tt.req)
			if recorder.Code != tt.code {
				t.Errorf("returned %v. Expected %v.", recorder.Code, tt.code)
			}
		})
	}
}

func TestPostValue(t *testing.T) {
	tests := []struct {
		name       string
		cacheStore *cacheStoreMock
		req        *http.Request
		code       int
	}{
		{
			name:       "should return 200",
			cacheStore: nil,
			req:        &http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/key"}, Body: io.NopCloser(strings.NewReader("välue"))},
			code:       http.StatusOK,
		},
		{
			name:       "should return 400",
			cacheStore: nil,
			req:        &http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/key"}, Body: nil},
			code:       http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		apiHandler := api.NewApiHandler(tt.cacheStore)
		recorder := httptest.NewRecorder()
		apiHandler.Router.ServeHTTP(recorder, tt.req)
		if recorder.Code != tt.code {
			t.Errorf("returned %v. Expected %v.", recorder.Code, tt.code)
		}
	}
}

func TestGetValue(t *testing.T) {
	tests := []struct {
		name       string
		cacheStore *cacheStoreMock
		req        *http.Request
		code       int
	}{
		{
			name:       "should return 200",
			cacheStore: nil,
			req:        &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/key1"}, Body: nil},
			code:       http.StatusOK,
		},
		{
			name: "should return 404",
			cacheStore: &cacheStoreMock{
				getFn: func(key string) (string, bool) {
					return "", false
				},
			},
			req:  &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/key2"}, Body: nil},
			code: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		apiHandler := api.NewApiHandler(tt.cacheStore)
		recorder := httptest.NewRecorder()
		apiHandler.Router.ServeHTTP(recorder, tt.req)
		if recorder.Code != tt.code {
			t.Errorf("returned %v. Expected %v.", recorder.Code, tt.code)
		}
	}
}
