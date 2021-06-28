package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/xid"
)

type requestIdKeyType string

const requestIdKey requestIdKeyType = "req-id"

func LogWithID(ctx context.Context, format string, args ...interface{}) {
	id, ok := ctx.Value(requestIdKey).(string)
	if !ok {
		log.Printf("failed to log request, no id")
		return
	}
	format = fmt.Sprintf("request %s: %s", id, format)
	log.Printf(format, args...)
}

// AppRouteHandlers define the handlers for py-server using the given dependencies
type AppRouteHandlers struct {
	TokenChecker   TokenChecker
	UserSaveStorer UserSaveStorer
}

func (h AppRouteHandlers) GetHandler(w http.ResponseWriter, req *http.Request) {
	authenticateRequest(h.TokenChecker, fetchHandler(h.UserSaveStorer))(w, req)
}

func (h AppRouteHandlers) PostHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, "coming soon")
}

func (h AppRouteHandlers) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, "coming soon")
}

// RouterHandlers are the possible handlers for the Router
type RouterHandlers interface {
	GetHandler(w http.ResponseWriter, req *http.Request)
	PostHandler(w http.ResponseWriter, req *http.Request)
	DeleteHandler(w http.ResponseWriter, req *http.Request)
}

// Route takes a set of RouteHandlers and routes a request to the appropriate handler
func Route(handler RouterHandlers, allowedOrigin string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if path != "/v1/usersave" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "not found")
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

		switch req.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Token")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			LogWithID(req.Context(), "served CORS options")
		case http.MethodGet:
			handler.GetHandler(w, req)
		case http.MethodPost:
			handler.PostHandler(w, req)
		case http.MethodDelete:
			handler.DeleteHandler(w, req)
		default:
			LogWithID(req.Context(), "invalid method used")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "invalid method")
		}
	}
}

// Serve runs ListenAndServe in a new goroutine, sending errors into shutdown,
// and blocking until ctx finishes before shutting down the server gracefully.
// Automatically logs all requests
func Serve(ctx context.Context, addr string, shutdown chan error, handler http.HandlerFunc) {
	server := http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			id := xid.New()
			log.Printf("entry %s: %s -> %s: host %s user-agent %s", id, req.Method, req.URL, req.Host, req.UserAgent())

			ctx := context.WithValue(req.Context(), requestIdKey, id.String())
			handler(w, req.WithContext(ctx))

			log.Printf("exit %s: request finished", id)
		}),
	}

	go func() {
		err := server.ListenAndServe()
		shutdown <- err
	}()

	<-ctx.Done()
	shutdown <- server.Shutdown(context.Background())
}
