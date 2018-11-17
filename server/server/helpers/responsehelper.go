package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	ContentTypeJSON = "application/json"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func JSONResponse(w http.ResponseWriter, data interface{}, httpCode int) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Failed to marshal JSON response: %v", err)
		panic(err)
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(httpCode)
	w.Write(jsonBytes)
}
