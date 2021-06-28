package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"py-server/server"
	"py-server/storage"
	"py-server/token"
)

type opts struct {
	development bool
	useFS       bool
}

func getOpts() opts {
	devFlag := flag.Bool("development", false, "Runs the server in development mode")
	fsFlag := flag.Bool("fs", false, "Use the filesystem to store UserSaves")

	flag.Parse()
	return opts{
		development: *devFlag,
		useFS:       *fsFlag,
	}
}

func getServerAddr() string {
	addr, found := os.LookupEnv("PYSERVER_ADDR")
	if !found {
		return "0.0.0.0:5000"
	}
	return addr
}

func main() {
	ctx := context.Background()
	opts := getOpts()

	allowedOrigin := os.Getenv("PYSERVER_ALLOWED_ORIGIN")
	serverAddr := getServerAddr()

	checkerClientID, foundClientID := os.LookupEnv("PYSERVER_CLIENTID")
	if !foundClientID && !opts.development {
		log.Fatal("client ID must be provided if server is not in development mode")
	}

	// Choose a UserSaveStorer implementation
	var storer server.UserSaveStorer
	if opts.useFS {
		storerRoot := os.Getenv("PYSERVER_STORE_ROOT")
		if len(storerRoot) == 0 {
			log.Fatal("no store root provided")
		}
		s, err := storage.MakeFileSystemStorer(storerRoot)
		if err != nil {
			log.Fatalf("failed to make storer: %s", err)
		}
		storer = s
	} else {
		s, err := storage.MakeGoogleStorer(ctx)
		if err != nil {
			log.Fatalf("failed to make storer: %s", err)
		}
		storer = s
	}

	routeHandlers := server.AppRouteHandlers{
		UserSaveStorer: storer,
		TokenChecker:   token.MakeGoogleTokenChecker(checkerClientID),
	}
	shutdownServer := make(chan error)
	go server.Serve(ctx, serverAddr, shutdownServer, server.Route(routeHandlers, allowedOrigin))

	log.Println(fmt.Sprintf("server started on %s", serverAddr))
	err := <-shutdownServer
	log.Println("exiting")
	if err != nil {
		log.Fatal(err)
	}
}
