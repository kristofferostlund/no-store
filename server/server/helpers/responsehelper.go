package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	ContentTypeJSON = "application/json"
)

func JSONResponse(w http.ResponseWriter, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Failed to marshal JSON response: %v", err)
		panic(err)
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	w.Write(jsonBytes)

}
