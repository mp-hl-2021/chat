package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_postSignup(t *testing.T) {
	m := postSignupRequestModel{
		Login:    "alice",
		Password: "123",
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal("failed to marshal struct")
	}
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(b))
	resp := httptest.NewRecorder()

	NewRouter().ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Server MUST return %d status code, but %d given", http.StatusCreated, resp.Code)
	}
	// note: this case will fail after hardcoded account id removal
	location := resp.Header().Get("Location")
	if  location != "/accounts/1" {
		t.Errorf("Server MUST return %s Location header, but %s given", "/accounts/1", location)
	}
}

func Test_postSignin(t *testing.T) {
	m := postSignupRequestModel{
		Login:    "alice",
		Password: "123",
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal("failed to marshal struct")
	}
	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(b))
	resp := httptest.NewRecorder()

	NewRouter().ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Server MUST return %d status code, but %d given", http.StatusOK, resp.Code)
	}
}