package defaults

import "net/http"

const ContentType = "application/json; charset=utf8"

func Headers(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", ContentType)
}
