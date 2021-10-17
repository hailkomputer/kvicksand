package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter().StrictSlash(false),
	}
	s.Router.Use(loggingMiddleware)
	s.routes()
	return s
}

func (s *Server) ServeHttp(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.Router.Path("/").
		HandlerFunc(s.handleIndex())
	s.Router.Path("/status").
		HandlerFunc(s.handleStatus())
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello!")
	}
}

func (s *Server) handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Status")
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
