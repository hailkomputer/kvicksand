package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *ApiHandler) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello!")
	}
}

func (a *ApiHandler) handlePostValue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)[keyName]
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
		key := mux.Vars(r)[keyName]
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
