package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func FromJSONBody(readable io.Reader, out interface{}) error {
	body, err := ioutil.ReadAll(readable)
	if err != nil {
		return fmt.Errorf("Failed to read request body: %v", err)
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return fmt.Errorf("Failed to unmarshal request body: %v", err)
	}

	return nil
}
