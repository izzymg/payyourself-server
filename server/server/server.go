package server

import (
	"context"
	"fmt"
	"net/http"
)

type TokenChecker interface {
	TokenIsValid(ctx context.Context, token string) bool
}

type GetHandler struct {
	tokenChecker TokenChecker
}

func (h GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get("Token")

	if len(token) < 1 {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "no token provided")
		return
	}

	if !h.tokenChecker.TokenIsValid(req.Context(), token) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "token invalid")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello!")
}

func Serve(ctx context.Context, addr string, shutdown chan error) {
	getHandler := GetHandler{}

	server := http.Server{
		Addr:    addr,
		Handler: getHandler,
	}

	go func() {
		err := server.ListenAndServe()
		shutdown <- err
	}()

	<-ctx.Done()
	shutdown <- server.Shutdown(context.Background())
}
