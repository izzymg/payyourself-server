package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

type TokenChecker interface {
	TokenIsValid(ctx context.Context, token string) (string, bool)
}

type UserSaveStorer interface {
	Fetch(userID string) (io.ReadCloser, error)
}

type AuthenticatedRequest struct {
	req    *http.Request
	userID string
}

type AuthenticatedRequestHandler = func(w http.ResponseWriter, req *AuthenticatedRequest)

// Returns an HTTP handler which checks the request token against the provided
// TokenChecker, calling next if it is valid, rejecting the request if invalid.
func checkRequestToken(tokenChecker TokenChecker, next AuthenticatedRequestHandler) http.HandlerFunc {
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
		next(w, &AuthenticatedRequest{
			userID: userID,
			req:    req,
		})
	}
}

type UserSaveHandler struct {
	userSaveStorer UserSaveStorer
}

func (h UserSaveHandler) HandleFetch(w http.ResponseWriter, req *AuthenticatedRequest) {
	reader, err := h.userSaveStorer.Fetch(req.userID)
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

func (h UserSaveHandler) HandleSave(w http.ResponseWriter, req *AuthenticatedRequest) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "soon to come")
}

func (h UserSaveHandler) HandleDelete(w http.ResponseWriter, req *AuthenticatedRequest) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "soon to come")
}

func MakeUserSaveHandler(userSaveStorer UserSaveStorer) UserSaveHandler {
	return UserSaveHandler{
		userSaveStorer,
	}
}
