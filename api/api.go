package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hailkomputer/kvicksand/pkg/cache"
)

type ApiHandler struct {
	Router *mux.Router
	Cache  *cache.Cache
}

func NewApiHandler() *ApiHandler {
	a := &ApiHandler{
		Router: mux.NewRouter().StrictSlash(false),
		Cache:  cache.NewCache(),
	}

	a.Router.Use(loggingMiddleware, recoveryMiddleware)

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

// recoveryMiddleware creates a HTTP handler for recovering from a panic.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			rec := recover()
			if rec != nil {
				switch t := rec.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				msg := fmt.Sprintf("Recovered: %s, request: %s %s", err, r.Method, r.URL.Path)
				http.Error(w, msg, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
