package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testTokenChecker struct {
	TestTokenIsValid func(token string) bool
}

func (t testTokenChecker) TokenIsValid(token string) bool {
	return t.TestTokenIsValid(token)
}

func testGetHandler(expectCode int, req *http.Request, tokenChecker TokenChecker) func(t *testing.T) {
	handler := GetHandler{tokenChecker}
	return func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if code := rr.Result().StatusCode; code != expectCode {
			t.Errorf("expected status %d, got %d", expectCode, code)
		}
	}
}

func TestGetHandler(t *testing.T) {

	matchingToken := "abc"
	tokenChecker := testTokenChecker{
		TestTokenIsValid: func(token string) bool {
			if token == matchingToken {
				return true
			} else {
				return false
			}
		},
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("failed to create test request: %w", err)
	}
	t.Run("no token", testGetHandler(http.StatusForbidden, req, tokenChecker))

	req.Header.Set("Token", "hello")
	t.Run("unmatching token", testGetHandler(http.StatusUnauthorized, req, tokenChecker))

	req.Header.Set("Token", matchingToken)
	t.Run("matching token", testGetHandler(http.StatusOK, req, tokenChecker))
}

func TestServe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go Serve(ctx, "localhost:5000")

	_, err := http.Get("http://localhost:5000")
	if err != nil {
		t.Error(err)
	}

	cancel()
}
