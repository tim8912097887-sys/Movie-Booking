package json

import (
	"encoding/json"
	"net/http"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/response"
)

func SuccessJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	payload, err := json.Marshal(response.NewSuccessResponse(data))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return err
	}
	w.WriteHeader(status)
	w.Write(payload)
	return nil
}