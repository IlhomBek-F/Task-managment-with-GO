package handler

import (
	"encoding/json"
	"net/http"
)

func render(w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	switch v := body.(type) {
	case string:
		json.NewEncoder(w).Encode(struct {
			Message string
		}{Message: v})
	case error:
		json.NewEncoder(w).Encode(struct {
			Error string
		}{
			Error: v.Error(),
		})

	case nil:
	default:
		json.NewEncoder(w).Encode(body)
	}
}
