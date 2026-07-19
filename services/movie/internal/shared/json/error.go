package json

import (
	"encoding/json"
	"net/http"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/response"
)

func ErrorJson(w http.ResponseWriter, status int, code string, message string) error {
	w.Header().Set("Content-Type", "application/json")
	payload, err := json.Marshal(response.NewErrorResponse(code, message))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return err
	}
	w.WriteHeader(status)
	w.Write(payload)
	return nil
}