package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testGetUser(t testing.T) {
	t.Run("should not allow unauthenticated  request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/user/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		
	})
}
