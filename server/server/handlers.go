package server

import (
	"context"
	"fmt"
	"net/http"
)

type TokenChecker interface {
	TokenIsValid(ctx context.Context, token string) bool
}

// Returns an HTTP handler which checks the request token against the provided
// TokenChecker, calling next if it is valid, rejecting the request if invalid.
func checkRequestToken(tokenChecker TokenChecker, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Token")

		if len(token) < 1 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "no token provided")
			return
		}

		if !tokenChecker.TokenIsValid(req.Context(), token) {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "token invalid")
			return
		}

		// valid token
		next(w, req)
	}
}

// PYHandler is the AppHandler for py-server
type PYHandler struct {
	tokenChecker TokenChecker
}

func (h PYHandler) GetHandler(w http.ResponseWriter, req *http.Request) {
	checkRequestToken(h.tokenChecker, func(rw http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintln(w, "soon to come")
	})
}

func (h PYHandler) PostHandler(w http.ResponseWriter, req *http.Request) {
	checkRequestToken(h.tokenChecker, func(rw http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintln(w, "soon to come")
	})
}

func (h PYHandler) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	checkRequestToken(h.tokenChecker, func(rw http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintln(w, "soon to come")
	})
}
