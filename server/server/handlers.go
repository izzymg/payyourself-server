package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"py-server/usersave"
)

var ErrNoUserSave = errors.New("no such user save")

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
	// Save returns a writer for the UserSave data at a given UserID
	// Writes should overwrite or create.
	Save(ctx context.Context, userID string) (io.WriteCloser, error)
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
		LogWithID(req.Context(), "trying to validate token")
		token := req.Header.Get("Token")

		if len(token) < 1 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "no token provided")
			LogWithID(req.Context(), "no token provided")
			return
		}

		userID, ok := tokenChecker.TokenIsValid(req.Context(), token)
		if !ok || len(userID) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "token invalid")
			LogWithID(req.Context(), "token invalid")
			return
		}

		if ok && len(userID) < 1 {
			LogWithID(req.Context(), "!! token validator returned ok, but user id is blank")
		}

		LogWithID(req.Context(), "validated token")
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
		LogWithID(req.req.Context(), "trying to fetch usersave")

		reader, err := userSaveStorer.Fetch(req.userID)
		if err != nil {
			if errors.Is(err, ErrNoUserSave) {
				LogWithID(req.req.Context(), "no usersave")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "No UserSave for this user")
				return
			}
			LogWithID(req.req.Context(), "!! failed to fetch usersave: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to fetch user save")
			return
		}
		defer func() {
			err := reader.Close()
			if err != nil {
				LogWithID(req.req.Context(), "!! failed to close usersave reader: %s", err)
			}
		}()

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, reader)
		if err != nil {
			LogWithID(req.req.Context(), "!! failed to send usersave: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to send user save")
			return
		}
		LogWithID(req.req.Context(), "sent usersave")
	}
}

// fetchHandler generates an AuthenticatedRequestHandler for saving with a
// UserSaveStorer
func saveHandler(userSaveStorer UserSaveStorer) authenticatedRequestHandler {
	return func(w http.ResponseWriter, req *authenticatedRequest) {
		LogWithID(req.req.Context(), "trying to save usersave")

		if req.req.Body == nil {
			LogWithID(req.req.Context(), "no body given")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "No body")
			return
		}

		// usersave is decoded from request body to validate correct schema
		userSave, err := usersave.DecodeUserSave(req.req.Body)
		if err != nil {
			LogWithID(req.req.Context(), "failed to decode incoming usersave: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Failed to decode user save")
			return
		}

		// usersave is re-encoded into the UserSaveStorer
		writer, err := userSaveStorer.Save(req.req.Context(), req.userID)
		if err != nil {
			if err != nil {
				LogWithID(req.req.Context(), "!! failed to get usersave writer: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Failed to send user save")
				return
			}
		}
		defer func() {
			err := writer.Close()
			if err != nil {
				LogWithID(req.req.Context(), "!! failed to close usersave writer: %s", err)
			}
		}()

		w.Header().Add("Content-Type", "application/json")
		err = usersave.EncodeUserSave(userSave, writer)
		if err != nil {
			LogWithID(req.req.Context(), "!! failed to encode outgoing usersave: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to send user save")
			return
		}
		LogWithID(req.req.Context(), "saved usersave")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Saved UserSave")
	}
}
