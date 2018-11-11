package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kristofferostlund/no-store/nostore"
	"github.com/kristofferostlund/no-store/server/server/helpers"

	"github.com/sirupsen/logrus"
)

var DefaultInterval = time.Duration(time.Minute * 10)

func (s *server) handleGetSecret() http.HandlerFunc {
	type response struct{}

	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		secrets, exists := query["secret"]
		if !exists {
			logrus.Warnf("Missing secrets parameter: %+v", query)
			http.Error(w, "secret is a required parameter", http.StatusBadRequest)
			return
		}

		decodedBytes, expired, err := nostore.Decode(secrets[0])
		if err != nil {
			logrus.Errorf("Failed to decode secret: %v", err)
			http.Error(w, "Can't decode secret", http.StatusBadRequest)
			return
		}

		if expired {
			http.Error(w, "Secret is expired", 419)
			return
		}

		w.Write(decodedBytes)
	}
}

func (s *server) handlePostSecret() http.HandlerFunc {
	type response struct {
		URL string `json:"url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ttl := DefaultInterval
		query := r.URL.Query()

		if ttlStrings, exists := query["ttl"]; exists {
			ttlSeconds, err := strconv.Atoi(ttlStrings[0])
			if err != nil {
				logrus.Warnf("Failed to parse ttl %s because %v", ttlStrings[0], err)
			}

			ttl = time.Duration(ttlSeconds) * time.Second
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Errorf("Failed to read body: %v", err)
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}

		compressed, err := nostore.Encode(body, time.Now().Add(ttl))
		if err != nil {
			logrus.Errorf("Failed to encode body: %v", err)
			http.Error(w, "Can't encode body", http.StatusBadRequest)
			return
		}

		url := fmt.Sprintf(
			"%s%s?secret=%s",
			r.Host,
			NoStoreSecretsPath,
			url.QueryEscape(string(compressed)),
		)

		helpers.JSONResponse(w, response{URL: url})
	}
}
