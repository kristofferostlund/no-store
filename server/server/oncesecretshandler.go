package server

import (
	"net/http"
	"time"
)

func (s *server) handleGetSecretOnce() http.HandlerFunc {
	type response struct {
		Value string `json:"value"`
	}

	type usedResponse struct {
		Reason string    `json:"reason"`
		UsedAt time.Time `json:"usedAt"`
	}

	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s *server) handlePostSecretOnce() http.HandlerFunc {
	type request struct {
		Value string `json:"value"`
	}

	type response struct {
		URL string `json:"url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {}
}
