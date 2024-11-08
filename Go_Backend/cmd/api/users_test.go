package main

import (
	"testing"
	"net/http"
)

func TestGetUser(t *testing.T) {
	app := newTestApp(t)
	mux := app.mount()
	
	testToken, err := app.authenticator.GenerateToken(nil) 
	if err != nil {
		t.Fail()
	}

	t.Run("should not allow unauthenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)	
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)	
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
		})
	t.Run("should allow authenticated requests", func (t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		
		req.Header.Set("Authorization", "Bearer " + testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
