package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func StatusCodeTest(req *http.Request, expect int, handler http.HandlerFunc) func(t *testing.T) {
	return func(t *testing.T) {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if code := rr.Code; code != expect {
			t.Errorf("expected status code %d, got %d", expect, code)
		}
	}
}

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

// returns status teapot to stub app handler for route testing
type teapotHandler struct{}

func (h teapotHandler) GetHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}
func (h teapotHandler) PostHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}
func (h teapotHandler) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}

// test that all routes except the API routes return not found
func TestRouter(t *testing.T) {

	handler := teapotHandler{}
	router := Router(handler)

	tests := map[string]int{
		"https://example.com/invalid":       http.StatusNotFound,
		"https://example.com/invalid/dog":   http.StatusNotFound,
		"https://example.com":               http.StatusNotFound,
		"abc":                               http.StatusNotFound,
		"https://example.com/":              http.StatusNotFound,
		"https://example.com/v1/usersave":   http.StatusTeapot,
		"http://localhost:1337/v1/usersave": http.StatusTeapot,
		"https://example.com/v1/abc":        http.StatusNotFound,
	}

	for route, expectedCode := range tests {
		t.Run(route, func(t *testing.T) {
			req, err := http.NewRequest("GET", route, nil)
			if err != nil {
				t.Error(err)
			}
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if code := rr.Code; code != expectedCode {
				t.Errorf("expected status code %d, got %d", expectedCode, rr.Code)
			}
		})
	}

	handler = teapotHandler{}
	router = Router(handler)

	allowedMethods := []string{http.MethodGet, http.MethodDelete, http.MethodPost}
	notAllowedMethods := []string{http.MethodConnect, http.MethodPut, http.MethodPatch}

	for _, method := range allowedMethods {
		req, err := http.NewRequest(method, "https://example.com/v1/usersave", nil)
		if err != nil {
			t.Error(err)
		}
		t.Run(method, StatusCodeTest(req, http.StatusTeapot, router))
	}

	for _, method := range notAllowedMethods {
		req, err := http.NewRequest(method, "https://example.com/v1/usersave", nil)
		if err != nil {
			t.Error(err)
		}
		t.Run(method, StatusCodeTest(req, http.StatusMethodNotAllowed, router))
	}
}

// test the server cycles up and down correctly
func TestServe(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	t.Run("should die on invalid port", func(t *testing.T) {
		shutdown := make(chan error)
		go Serve(context.TODO(), "localhost:-1", shutdown, http.HandlerFunc(handler))
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

			go Serve(ctx, "localhost:0", shutdown, http.HandlerFunc(handler))
			cancel()
			err := <-shutdown
			if err != nil {
				t.Error("got err on graceful shutdown: %w", err)
			}
		})
	}
}
