package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/greg-szabo/f11/defaults"
	"net/http"
)

type ClaimMessageV1 struct {
	Status string `json:"status"`
}

func ClaimHandlerV1(w http.ResponseWriter, r *http.Request) {
	//Todo: Implement gaiacli tx send
	defaults.Headers(w)
	json.NewEncoder(w).Encode(ClaimMessageV1{"submitted"})
}

func AddRoutesV1(r *mux.Router) {

	r.HandleFunc("/v1/claim", ClaimHandlerV1)

}
