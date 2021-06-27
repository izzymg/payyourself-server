package server

import (
	"context"
	"net/http"
	"testing"
)

type testTokenChecker struct {
	TestTokenIsValid func(ctx context.Context, token string) bool
}

func (t testTokenChecker) TokenIsValid(ctx context.Context, token string) bool {
	return t.TestTokenIsValid(ctx, token)
}

// checkRequestToken with next handler that returns status teapot
func makeCheckTokenTest(expectCode int, req *http.Request, tokenChecker TokenChecker) func(t *testing.T) {
	return func(t *testing.T) {
		handler := checkRequestToken(tokenChecker, func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		})

		t.Run("check request token", StatusCodeTest(req, expectCode, handler))
	}
}

// test that CheckRequestToken doesn't allow requests to pass through
// unless the provided TokenChecker approves it
func TestCheckRequestToken(t *testing.T) {
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
	t.Run("no token", makeCheckTokenTest(http.StatusForbidden, req, tokenChecker))

	req.Header.Set("Token", "hello")
	t.Run("unmatching token", makeCheckTokenTest(http.StatusUnauthorized, req, tokenChecker))

	req.Header.Set("Token", matchingToken)
	t.Run("matching token", makeCheckTokenTest(http.StatusTeapot, req, tokenChecker))
}

func TestPYHandler(t *testing.T) {
	// token checker that always allows requests
	tokenChecker := testTokenChecker{
		TestTokenIsValid: func(ctx context.Context, token string) bool {
			return true
		},
	}

	pyHandler := MakePYHandler(tokenChecker)

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Error(err)
	}

	StatusCodeTest(
		req,
		http.StatusNotImplemented,
		pyHandler.GetHandler,
	)
	StatusCodeTest(
		req,
		http.StatusNotImplemented,
		pyHandler.PostHandler,
	)
	StatusCodeTest(
		req,
		http.StatusNotImplemented,
		pyHandler.DeleteHandler,
	)
}
