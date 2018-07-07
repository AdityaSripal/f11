package main

import (
	"github.com/greg-szabo/f11/testnet"
	"github.com/greg-szabo/f11/version"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_MainHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "{\"name\":\"" + testnet.TestnetName + "\",\"version\":\"" + version.Version + "\"}\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
