package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testTokenChecker struct {
	TestTokenIsValid func(ctx context.Context, token string) (string, bool)
}

func (t testTokenChecker) TokenIsValid(ctx context.Context, token string) (string, bool) {
	return t.TestTokenIsValid(ctx, token)
}

// checkRequestToken with next handler that returns status teapot
func makeCheckTokenTest(expectCode int, req *http.Request, tokenChecker TokenChecker) func(t *testing.T) {
	return func(t *testing.T) {
		handler := authenticateRequest(tokenChecker, func(w http.ResponseWriter, req *authenticatedRequest) {
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
		TestTokenIsValid: func(ctx context.Context, token string) (string, bool) {
			if token == matchingToken {
				return "someID", true
			} else {
				return "", false
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

// fake io.readcloser
type testReadCloser struct{}

func (t testReadCloser) Read(p []byte) (int, error) {
	return 0, io.EOF
}
func (t testReadCloser) Close() error {
	return nil
}

// fake storer of user saves that fetches the fake readercloser
type testUserSaveStorer struct{}

func (t testUserSaveStorer) Fetch(userID string) (io.ReadCloser, error) {
	return testReadCloser{}, nil
}

func TestHandleFetch(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	authedReq := authenticatedRequest{
		req:    req,
		userID: "some user id",
	}

	rr := httptest.NewRecorder()
	fetchHandler(testUserSaveStorer{})(rr, &authedReq)

	if code := rr.Code; code != http.StatusOK {
		t.Errorf("expected code %d, got %d", http.StatusOK, code)
	}
}
