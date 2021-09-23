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
}

func getOpts() opts {
	devFlag := flag.Bool("development", false, "Runs the server in development mode")

	flag.Parse()
	return opts{
		development: *devFlag,
	}
}

func getServerAddr() string {
	addr, found := os.LookupEnv("PYSERVER_ADDR")
	if !found {
		return "0.0.0.0:5000"
	}
	return addr
}

func getBucketName() string {
	name, found := os.LookupEnv("PYSERVER_BUCKET_NAME")
	if !found {
		return "user-saves-1"
	}
	return name
}

func main() {
	ctx := context.Background()
	opts := getOpts()

	allowedOrigin := os.Getenv("PYSERVER_ALLOWED_ORIGIN")
	serverAddr := getServerAddr()
	bucketName := getBucketName()

	checkerClientID, foundClientID := os.LookupEnv("PYSERVER_CLIENTID")
	if !foundClientID && !opts.development {
		log.Fatal("client ID must be provided if server is not in development mode")
	}

	log.Println("bringing up google cloud storer")
	storer, err := storage.MakeGoogleStorer(ctx, bucketName)
	if err != nil {
		log.Fatalf("failed to make storer: %s", err)
	}
	defer storer.Close()
	log.Println("google cloud storer up")

	routeHandlers := server.AppRouteHandlers{
		UserSaveStorer: storer,
		TokenChecker:   token.MakeGoogleTokenChecker(checkerClientID),
	}
	shutdownServer := make(chan error)
	go server.Serve(ctx, serverAddr, shutdownServer, server.Route(routeHandlers, allowedOrigin))

	log.Println(fmt.Sprintf("server started on %s", serverAddr))
	err = <-shutdownServer
	log.Println("exiting")
	if err != nil {
		log.Fatal(err)
	}
}
