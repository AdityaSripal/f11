package defaults

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func Headers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", ContentType)
}

func AddStatusOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func InternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorMessage{err.Error()})
}

func UserError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorMessage{err.Error()})
}
