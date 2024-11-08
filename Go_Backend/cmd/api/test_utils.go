package main

import (
	"CoffeeMap/internal/store"
	"testing"
	"net/http/httptest"
	"net/http"
	"CoffeeMap/internal/auth"
)

func newTestApp(t *testing.T) *application {
	t.Helper()

	mockStore := store.NewMockStore()
	testAuth := &auth.TestAuthenticator{}
	return &application {
		store: mockStore,
		authenticator: testAuth,
	}
}


func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, got %d", expected, actual)
	}
}
