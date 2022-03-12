package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CacheStore interface {
	// Get tries to fetch the value for the given key
	// It will always return second argument as false, in cases where first
	// value is expired or not found
	Get(key string) (string, bool)
	// Set writes the value for the given key
	// Expiration duration is hard coded as 30 minutes
	// If a value for the specified key already exists, then it will be overwritten
	Set(key, value string)
}

type ApiHandler struct {
	Router *mux.Router
	Cache  CacheStore
}

func NewApiHandler(cache CacheStore) *ApiHandler {
	a := &ApiHandler{
		Router: mux.NewRouter().StrictSlash(false),
		Cache:  cache,
	}

	a.Router.Use(loggingMiddleware)

	for _, route := range a.createRoutes() {
		a.Router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
		log.Printf("registered endpoint %s:%s", route.Method, route.Name)
	}
	return a
}

func (a *ApiHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}

// loggingMiddleware creates a middleware for logging each request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s:%s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}
