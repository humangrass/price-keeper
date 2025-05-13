package xhttp

import (
	"net/http"

	"github.com/mailru/easyjson"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload easyjson.Marshaler) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err := easyjson.MarshalToWriter(payload, w)
	if err != nil {
		return err
	}

	return nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) error {
	return RespondWithJSON(w, code, &ErrorResponse{
		Error: message,
	})
}

//easyjson:json
type ErrorResponse struct {
	Error string `json:"error"`
}
