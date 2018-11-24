package server

import (
	"fmt"
	"net/http"
)

const BasePath = "/api/v0"

var (
	HealthPath      = BasePath + "/health"
	TTLSecretsPath  = BasePath + "/secrets/ttl"
	OnceSecretsPath = BasePath + "/secrets/once"
)

func (s *server) initRoutes() {
	s.router.HandleFunc("GET", HealthPath, s.handleHealth())
	s.router.HandleFunc("GET", TTLSecretsPath, s.handleGetSecret())
	s.router.HandleFunc("POST", TTLSecretsPath, s.handlePostSecret())
	s.router.HandleFunc("GET", OnceSecretsPath, s.handleGetSecretOnce())
	s.router.HandleFunc("POST", OnceSecretsPath, s.handlePostSecretOnce())

	s.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not found\n")
	})
}
