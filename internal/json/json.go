// Package json provides helper functions for JSON encoding and decoding.
package json

import (
	"encoding/json"
	"net/http"
)

// Write sends a JSON response with the given status code and data.
func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Read decodes the JSON body of the request into the given data structure.
func Read(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}
