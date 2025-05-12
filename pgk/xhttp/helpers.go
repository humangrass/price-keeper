package xhttp

import (
	"net/http"

	"github.com/mailru/easyjson"
)

// func RespondWithJSON(w http.ResponseWriter, code int, payload any) error {
// 	response, err := json.Marshal(payload)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return err
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	_, err = w.Write(response)
// 	return err
// }

func RespondWithJSON(w http.ResponseWriter, code int, payload easyjson.Marshaler) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err := easyjson.MarshalToWriter(payload, w)
	if err != nil {
		// Если ошибка произошла после отправки заголовков,
		// мы не можем изменить статус, только логировать
		return err
	}

	return nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) error {
	return RespondWithJSON(w, code, &ErrorResponse{
		Error: message,
	})
}

// func RespondWithError(w http.ResponseWriter, code int, message string) error {
// 	return RespondWithJSON(w, code, map[string]string{"error": message})
// }

//easyjson:json
type ErrorResponse struct {
	Error string `json:"error"`
}
