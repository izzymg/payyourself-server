package server

import (
	"context"
	"fmt"
	"net/http"
)

// AppHandler defines the handling functions of py-server.
type AppHandler interface {
	GetHandler(w http.ResponseWriter, req *http.Request)
	PostHandler(w http.ResponseWriter, req *http.Request)
	DeleteHandler(w http.ResponseWriter, req *http.Request)
}

// Router is the main route handler for py-server.
func Router(handler AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if path != "/v1/usersave" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "not found")
			return
		}

		switch req.Method {
		case http.MethodGet:
			handler.GetHandler(w, req)
		case http.MethodPost:
			handler.PostHandler(w, req)
		case http.MethodDelete:
			handler.DeleteHandler(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "invalid method")
			return
		}
	}
}

// Serve runs ListenAndServe in a new goroutine, sending errors into shutdown,
// and blocking until ctx finishes before shutting down the server gracefully.
func Serve(ctx context.Context, addr string, shutdown chan error, handler http.Handler) {
	server := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		err := server.ListenAndServe()
		shutdown <- err
	}()

	<-ctx.Done()
	shutdown <- server.Shutdown(context.Background())
}
