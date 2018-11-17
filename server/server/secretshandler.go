package server

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/kristofferostlund/no-store/nostore"
	"github.com/kristofferostlund/no-store/server/server/helpers"

	"github.com/sirupsen/logrus"
)

var DefaultInterval = time.Duration(time.Minute * 10)

func (s *server) handleGetSecret() http.HandlerFunc {
	type response struct {
		Value     string    `json:"value"`
		ExpiresAt time.Time `json:"expiresAt"`
		ExpiresIn int64     `json:"expiresIn"`
	}

	type expiredResponse struct {
		Reason    string    `json:"reason"`
		ExpiredAt time.Time `json:"expiredAt"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		secrets, exists := query["secret"]
		if !exists {
			logrus.Warnf("Missing secrets parameter: %+v", query)
			helpers.ErrorJSONResponse(w, "secret is a required parameter", http.StatusBadRequest)
			return
		}

		decodedBytes, expiresAt, expired, err := nostore.Decode(secrets[0])
		if err != nil {
			logrus.Errorf("Failed to decode secret: %v", err)
			helpers.ErrorJSONResponse(w, "Can't decode secret", http.StatusBadRequest)
			return
		}

		if expired {
			re := expiredResponse{
				Reason:    "Secret is expired",
				ExpiredAt: expiresAt,
			}

			helpers.JSONResponse(w, re, http.StatusGone)
			return
		}

		helpers.JSONResponse(
			w,
			response{
				Value:     string(decodedBytes),
				ExpiresAt: expiresAt,
				ExpiresIn: int64(math.Abs(time.Now().Sub(expiresAt).Seconds())),
			},
			http.StatusOK,
		)
	}
}

func (s *server) handlePostSecret() http.HandlerFunc {
	type request struct {
		Value      string `json:"value"`
		TTLSeconds int64  `json:"ttlSeconds"`
	}

	type response struct {
		URL       string    `json:"url"`
		ExpiresAt time.Time `json:"expiresAt"`
		ExpiresIn int64     `json:"expiresIn"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		inbound := request{}
		ttl := DefaultInterval

		if err := helpers.FromJSONBody(r.Body, &inbound); err != nil {
			logrus.Error(err)
			helpers.ErrorJSONResponse(w, "Can't read body, is the format correct?", http.StatusBadRequest)
			return
		}

		if inbound.TTLSeconds != 0 {
			ttl = time.Duration(inbound.TTLSeconds) * time.Second
		}

		resp := response{
			ExpiresAt: time.Now().Add(ttl),
			ExpiresIn: int64(ttl.Seconds()),
		}

		compressed, err := nostore.Encode([]byte(inbound.Value), resp.ExpiresAt)
		if err != nil {
			logrus.Errorf("Failed to encode body: %v", err)
			helpers.ErrorJSONResponse(w, "Can't encode body", http.StatusBadRequest)
			return
		}

		resp.URL = fmt.Sprintf(
			"%s%s?secret=%s",
			r.Host,
			NoStoreSecretsPath,
			url.QueryEscape(string(compressed)),
		)

		helpers.JSONResponse(w, resp, http.StatusOK)
	}
}
