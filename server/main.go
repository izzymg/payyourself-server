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

func isDevMode() bool {
	devFlag := flag.Bool("development", false, "Runs the server in development mode")
	flag.Parse()
	return *devFlag
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
	shutdownServer := make(chan error)

	allowedOrigin := os.Getenv("PYSERVER_ALLOWED_ORIGIN")

	storerRoot := os.Getenv("PYSERVER_STORE_ROOT")
	if len(storerRoot) == 0 {
		log.Fatal("no store root provided")
	}

	serverAddr := getServerAddr()
	checkerClientID, foundClientID := os.LookupEnv("PYSERVER_CLIENTID")
	if !foundClientID && !isDevMode() {
		log.Fatal("client ID must be provided if server is not in development mode")
	}

	storer, err := storage.MakeFileSystemStorer(storerRoot)
	if err != nil {
		log.Fatalf("failed to make storer: %s", err)
	}

	routeHandlers := server.AppRouteHandlers{
		UserSaveStorer: storer,
		TokenChecker:   token.MakeGoogleTokenChecker(checkerClientID),
	}
	go server.Serve(ctx, serverAddr, shutdownServer, server.Route(routeHandlers, allowedOrigin))

	log.Println(fmt.Sprintf("server started on %s", serverAddr))
	err = <-shutdownServer
	log.Println("exiting")
	if err != nil {
		log.Fatal(err)
	}
}
