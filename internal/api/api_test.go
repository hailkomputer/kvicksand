package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hailkomputer/kvicksand/internal/api"
	"github.com/hailkomputer/kvicksand/internal/api/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ApiHandlerSuite struct {
	suite.Suite

	recorder   *httptest.ResponseRecorder
	apiHandler *api.ApiHandler
	cache      *mocks.Cache
}

func (a *ApiHandlerSuite) SetupTest() {
	a.recorder = httptest.NewRecorder()
	a.cache = &mocks.Cache{}
	a.apiHandler = api.NewApiHandler(a.cache)
}

func TestApiHandlerSuite(t *testing.T) {
	suite.Run(t, new(ApiHandlerSuite))
}

func (a *ApiHandlerSuite) TestIndex_ShouldReturn200() {
	req, _ := http.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	a.apiHandler.Router.ServeHTTP(a.recorder, req)
	assert.Equal(a.T(), http.StatusOK, a.recorder.Code)
}

func (a *ApiHandlerSuite) TestPostValue_ShouldReturn200() {
	req, _ := http.NewRequest(
		http.MethodPost,
		"/key",
		strings.NewReader("välue"),
	)

	a.cache.On("Set", "key", "välue")

	a.apiHandler.Router.ServeHTTP(a.recorder, req)
	assert.Equal(a.T(), http.StatusOK, a.recorder.Code)
}

func (a *ApiHandlerSuite) TestPostValue_ShouldReturn400() {
	req, _ := http.NewRequest(
		http.MethodPost,
		"/key",
		nil,
	)

	a.apiHandler.Router.ServeHTTP(a.recorder, req)
	assert.Equal(a.T(), http.StatusBadRequest, a.recorder.Code)
}

func (a *ApiHandlerSuite) TestGetValue_ShouldReturn404() {
	req, _ := http.NewRequest(
		http.MethodGet,
		"/key",
		nil,
	)

	a.cache.On("Get", "key").Return("", false)

	a.apiHandler.Router.ServeHTTP(a.recorder, req)
	assert.Equal(a.T(), http.StatusNotFound, a.recorder.Code)
}

func (a *ApiHandlerSuite) TestGetValue_ShouldReturn200() {
	req, _ := http.NewRequest(
		http.MethodPost,
		"/key",
		strings.NewReader("välue"),
	)

	a.cache.On("Set", "key", "välue")
	a.apiHandler.Router.ServeHTTP(a.recorder, req)

	req, _ = http.NewRequest(
		http.MethodGet,
		"/key",
		nil,
	)

	a.cache.On("Get", "key").Return("välue", true)

	a.apiHandler.Router.ServeHTTP(a.recorder, req)
	assert.Equal(a.T(), http.StatusOK, a.recorder.Code)
}
