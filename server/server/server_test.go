package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testTokenChecker struct {
	TestTokenIsValid func(ctx context.Context, token string) bool
}

func (t testTokenChecker) TokenIsValid(ctx context.Context, token string) bool {
	return t.TestTokenIsValid(ctx, token)
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
		TestTokenIsValid: func(ctx context.Context, token string) bool {
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
	t.Run("should die on invalid port", func(t *testing.T) {
		shutdown := make(chan error)
		go Serve(context.TODO(), "localhost:-1", shutdown)
		err := <-shutdown
		if err == nil {
			t.Error("expected invalid port error, got none")
		}
	})

	testCount := 10
	for i := 0; i < testCount; i++ {
		t.Run(fmt.Sprintf("should startup,shutdown on cancel: %d", i), func(t *testing.T) {
			shutdown := make(chan error)
			ctx, cancel := context.WithCancel(context.Background())

			go Serve(ctx, "localhost:0", shutdown)
			cancel()
			err := <-shutdown
			if err != nil {
				t.Error("got err on graceful shutdown: %w", err)
			}
		})
	}
}
