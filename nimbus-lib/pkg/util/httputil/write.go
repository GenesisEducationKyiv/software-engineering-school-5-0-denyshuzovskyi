package httputil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	if err := enc.Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(buf.Bytes()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
