package api

import (
	"fmt"
	"net/http"
)

// Route holds API routing information.
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

// Routes is a list of routes.
type Routes []Route

const (
	keyName = "key"
)

func (a ApiHandler) createRoutes() Routes {
	return Routes{
		Route{
			Name:    "Index",
			Method:  http.MethodGet,
			Pattern: "/",
			Handler: http.HandlerFunc(a.handleIndex()),
		},
		Route{
			Name:    "PostValue",
			Method:  http.MethodPost,
			Pattern: fmt.Sprintf("/{%s}", keyName),
			Handler: http.HandlerFunc(a.handlePostValue()),
		},
		Route{
			Name:    "GetValue",
			Method:  http.MethodGet,
			Pattern: fmt.Sprintf("/{%s}", keyName),
			Handler: http.HandlerFunc(a.handleGetValue()),
		},
	}
}
