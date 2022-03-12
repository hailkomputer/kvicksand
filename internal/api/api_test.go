package api_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hailkomputer/kvicksand/internal/api"
	"github.com/hailkomputer/kvicksand/internal/api/mocks"
)

var (
	apiHandler *api.ApiHandler
	cache      *mocks.CacheStore
)

func TestMain(m *testing.M) {
	cache = &mocks.CacheStore{}
	apiHandler = api.NewApiHandler(cache)

	os.Exit(m.Run())
}

func TestIndex(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	apiHandler.Router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("returned %v. Expected %v.", recorder.Code, http.StatusOK)
	}
}

func TestPostValue(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/key", strings.NewReader("välue"))
	cache.On("Set", "key", "välue")
	recorder := httptest.NewRecorder()

	apiHandler.Router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("returned %v. Expected %v.", recorder.Code, http.StatusOK)
	}

	req, _ = http.NewRequest(http.MethodPost, "/key", nil)
	recorder = httptest.NewRecorder()

	apiHandler.Router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("returned %v. Expected %v.", recorder.Code, http.StatusBadRequest)
	}
}

func TestGetValue(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/key1", nil)
	cache.On("Get", "key1").Return("välue", true)
	recorder := httptest.NewRecorder()

	apiHandler.Router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Errorf("returned %v. Expected %v.", recorder.Code, http.StatusOK)
	}

	req, _ = http.NewRequest(http.MethodGet, "/key2", nil)
	cache.On("Get", "key2").Return("", false)
	recorder = httptest.NewRecorder()

	apiHandler.Router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("returned %v. Expected %v.", recorder.Code, http.StatusNotFound)
	}
}
