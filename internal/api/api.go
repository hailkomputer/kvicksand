package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// CacheStore is an interface for the cache store
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

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type routes []route

func (a *ApiHandler) createRoutes() routes {
	return routes{
		route{
			method:  http.MethodGet,
			regex:   regexp.MustCompile(`^/$`),
			handler: a.handleIndex(),
		},
		route{
			method:  http.MethodPost,
			regex:   regexp.MustCompile(`^/(?P<key>[a-zA-Z0-9]+)$`),
			handler: a.handlePostValue(),
		},
		route{
			method:  http.MethodGet,
			regex:   regexp.MustCompile(`^/(?P<key>[a-zA-Z0-9]+)$`),
			handler: a.handleGetValue(),
		},
	}
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

// ApiHandler is the main handler for the API
type ApiHandler struct {
	Cache  CacheStore
	Router http.Handler
}

// NewApiHandler creates a new ApiHandler
func NewApiHandler(cache CacheStore) *ApiHandler {
	a := &ApiHandler{
		Cache: cache,
	}

	a.Router = loggingMiddleware(http.HandlerFunc(a.Serve))
	return a
}

// Serve is the main router for the api
func (a *ApiHandler) Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range a.createRoutes() {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%s:%s", r.Method, r.RequestURI))
		next.ServeHTTP(w, r)
	})
}

func (a *ApiHandler) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello!")
	}
}

func (a *ApiHandler) handlePostValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := getField(r, 0)
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		a.Cache.Set(key, string(body))
	}
}

func (a *ApiHandler) handleGetValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := getField(r, 0)
		value, inCache := a.Cache.Get(key)
		switch inCache {
		case false:
			w.WriteHeader(http.StatusNotFound)
			return
		case true:
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			if _, err := io.WriteString(w, value); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
