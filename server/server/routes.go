package server

import (
	"fmt"
	"net/http"
)

const BasePath = "/api/v0"

var (
	HealthPath         = BasePath + "/health"
	NoStoreSecretsPath = BasePath + "/secrets"
)

func (s *server) initRoutes() {
	s.router.HandleFunc("GET", HealthPath, s.handleHealth())
	s.router.HandleFunc("GET", NoStoreSecretsPath, s.handleGetSecret())
	s.router.HandleFunc("POST", NoStoreSecretsPath, s.handlePostSecret())

	s.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not found\n")
	})
}
