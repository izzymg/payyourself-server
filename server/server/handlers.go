package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

// TokenChecker defines methods for validating a given token,
// providing the associated ID and true if valid, or false if invalid
type TokenChecker interface {
	TokenIsValid(ctx context.Context, token string) (string, bool)
}

// UserSaveStorer defines methods for fetching, deleting and saving
// UserSave data
type UserSaveStorer interface {
	// Fetch returns a reader for the UserSave data at a given UserID
	Fetch(userID string) (io.ReadCloser, error)
}

// authenticatedRequest wraps an HTTP request with a UserID
type authenticatedRequest struct {
	req    *http.Request
	userID string
}

type authenticatedRequestHandler = func(w http.ResponseWriter, req *authenticatedRequest)

// Returns an HTTP handler which checks the request token against the provided
// TokenChecker, calling next if it is valid, rejecting the request if invalid.
func authenticateRequest(tokenChecker TokenChecker, next authenticatedRequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Token")

		if len(token) < 1 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "no token provided")
			return
		}

		userID, ok := tokenChecker.TokenIsValid(req.Context(), token)
		if !ok || len(userID) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "token invalid")
			if len(userID) < 1 {
				log.Printf("token validated with empty user ID: %v", req)
			}
			return
		}

		// valid token
		next(w, &authenticatedRequest{
			userID: userID,
			req:    req,
		})
	}
}

// fetchHandler generates an AuthenticatedRequestHandler for fetching from the
// UserSaveStorer
func fetchHandler(userSaveStorer UserSaveStorer) authenticatedRequestHandler {
	return func(w http.ResponseWriter, req *authenticatedRequest) {
		reader, err := userSaveStorer.Fetch(req.userID)
		defer func() {
			err := reader.Close()
			if err != nil {
				log.Printf("failed to close user save: %s", err)
			}
		}()
		if err != nil {
			log.Printf("failed to fetch user save: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to fetch user save")
			return
		}

		_, err = io.Copy(w, reader)
		if err != nil {
			log.Printf("failed to send user save: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to send user save")
			return
		}
	}
}
